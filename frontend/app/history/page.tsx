"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Camera, Download, TrendingUp, Calendar, Clock, FileImage, FileVideo, ArrowLeft, Trash, Eye } from "lucide-react"
import Link from "next/link"
import { useAuth } from "@/components/AuthContext"
import { LoginModal } from "@/components/login-modal"
import { RegistrationModal } from "@/components/registration-modal"

const API_BASE = "http://localhost:8080"

export default function HistoryPage() {
  const { token, isAuthenticated, user, logout } = useAuth()
  const [files, setFiles] = useState<any[]>([])
  const [stats, setStats] = useState<{ total_files: number; processed_this_week: number }>({
    total_files: 0,
    processed_this_week: 0,
  })
  const [loading, setLoading] = useState(true)
  const [showLoginModal, setShowLoginModal] = useState(false)
  const [showRegistrationModal, setShowRegistrationModal] = useState(false)

  // Загрузка списка файлов
  const fetchFiles = async () => {
    if (!isAuthenticated || !token) {
      setLoading(false)
      return
    }
    try {
      setLoading(true)
      const res = await fetch(`${API_BASE}/api/files`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) {
        if (res.status === 401) {
          alert("Сессия истекла, пожалуйста, войдите снова")
          logout()
        }
        throw new Error(`Ошибка получения файлов: ${res.status}`)
      }
      const { data } = await res.json()
      setFiles(data || [])
    } catch (error) {
      console.error("Ошибка загрузки файлов:", error)
      alert("Не удалось загрузить историю файлов")
    } finally {
      setLoading(false)
    }
  }

  // Загрузка статистики
  const fetchStats = async () => {
    if (!isAuthenticated || !token) return
    try {
      const res = await fetch(`${API_BASE}/api/user/stats`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) {
        if (res.status === 401) {
          alert("Сессия истекла, пожалуйста, войдите снова")
          logout()
        }
        throw new Error(`Ошибка получения статистики: ${res.status}`)
      }
      const { data } = await res.json()
      setStats({
        total_files: data.total_files || 0,
        processed_this_week: data.processed_this_week || 0,
      })
    } catch (error) {
      console.error("Ошибка загрузки статистики:", error)
      setStats({ total_files: 0, processed_this_week: 0 })
    }
  }

  // Скачивание файла
  const handleDownload = async (fileId: string, type: "original" | "processed") => {
    try {
      const res = await fetch(`${API_BASE}/api/files/${fileId}?type=${type}`, {
        method: "GET",
        headers: isAuthenticated && token ? { Authorization: `Bearer ${token}` } : undefined,
      })

      if (!res.ok) {
        const error = await res.json().catch(() => ({}))
        if (res.status === 401 && isAuthenticated) {
          alert("Сессия истекла, пожалуйста, войдите снова")
          logout()
        }
        throw new Error(error.message || `Ошибка скачивания: ${res.status}`)
      }

      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement("a")
      a.href = url

      const file = files.find((f) => f.id === fileId)
      const extension = file?.original_name
        ? file.original_name.split(".").pop() || "file"
        : "file"
      a.download = `${type}-${fileId}.${extension}`
      a.click()
      URL.revokeObjectURL(url)
    } catch (error: any) {
      console.error(`Ошибка скачивания (${type}):`, error)
      alert(`Ошибка: ${error.message || "Не удалось скачать файл"}`)
    }
  }

  // Удаление файла
  const handleDelete = async (fileId: string) => {
    if (!isAuthenticated || !token) return
    try {
      const res = await fetch(`${API_BASE}/api/files/${fileId}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) {
        if (res.status === 401) {
          alert("Сессия истекла, пожалуйста, войдите снова")
          logout()
        }
        throw new Error(`Ошибка удаления: ${res.status}`)
      }
      setFiles(files.filter((f) => f.id !== fileId))
    } catch (error) {
      console.error("Ошибка удаления:", error)
      alert("Не удалось удалить файл")
    }
  }

  // Форматирование даты
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  // Получение типа файла
  const getFileType = (mimeType: string) => {
    return mimeType.startsWith("image/") ? "image" : mimeType.startsWith("video/") ? "video" : "file"
  }

  // Получение размера файла
  const formatFileSize = (size: number) => {
    if (size < 1024) return `${size} Б`
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} КБ`
    return `${(size / (1024 * 1024)).toFixed(1)} МБ`
  }

  // Маппинг статусов
  const statusMap: { [key: string]: string } = {
    uploaded: "Загружен",
    processing: "Обрабатывается",
    completed: "Завершено",
    failed: "Ошибка",
  }
/*
  // Маппинг типов размытия
  const blurTypeMap: { [key: string]: string } = {
    gaussian: "Размытие",
    pixelate: "Пикселизация",
    mask: "Маска",
  }

  // Маппинг качества
  const qualityMap: { [key: string]: string } = {
    original: "Оригинальное",
    high: "Высокое",
    medium: "Среднее",
  }
*/

  // Загрузка файлов и статистики при монтировании
  useEffect(() => {
    if (isAuthenticated && token) {
      fetchFiles()
      fetchStats()
    }
    const createParticles = () => {
      const particlesContainer = document.getElementById("particles")
      if (!particlesContainer) return

      for (let i = 0; i < 30; i++) {
        const particle = document.createElement("div")
        particle.className = "particle"
        particle.style.left = Math.random() * 100 + "%"
        particle.style.animationDelay = Math.random() * 15 + "s"
        particle.style.animationDuration = Math.random() * 10 + 10 + "s"
        particlesContainer.appendChild(particle)
      }
    }
    createParticles()
  }, [isAuthenticated, token])

  // Обработка переключения на регистрацию
  const handleSwitchToRegister = () => {
    setShowLoginModal(false)
    setTimeout(() => setShowRegistrationModal(true), 100)
  }

  // Обработка переключения на логин
  const handleSwitchToLogin = () => {
    setShowRegistrationModal(false)
    setTimeout(() => setShowLoginModal(true), 100)
  }

  return (
    <div className="min-h-screen bg-black relative">
      {/* Background Animation */}
      <div className="bg-animation absolute top-0 left-0 w-full h-full z-0"></div>
      <div className="particles absolute top-0 left-0 w-full overflow-hidden h-full z-10" id="particles"></div>

      {/* Header */}
      <header className="fixed top-0 w-full z-50 backdrop-blur-lg bg-black/10 border-b border-white/10 transition-all duration-300">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-4">
              <Link href="/" className="flex items-center gap-2 text-white/80 hover:text-white transition-colors">
                <ArrowLeft className="w-5 h-5" />
                <span className="font-manrope">Назад</span>
              </Link>
              <div className="w-px h-6 bg-white/20"></div>
              <div className="flex items-center gap-2">
                <div className="w-8 h-8 bg-gradient-to-br from-white to-purple-400 rounded-lg flex items-center justify-center">
                  <Camera className="w-5 h-5 text-black" />
                </div>
                <span className="font-geist font-bold text-xl bg-gradient-to-r from-white to-purple-400 bg-clip-text text-transparent">
                  Obscura
                </span>
              </div>
            </div>
            <div className="flex items-center gap-3">
              {isAuthenticated ? (
                <>
                  <span className="font-manrope text-white/80">
                    Привет, {user?.name || user?.email || "Пользователь"}!
                  </span>
                  <Button
                    variant="ghost"
                    onClick={logout}
                    className="font-manrope text-white hover:bg-white/10"
                  >
                    Выйти
                  </Button>
                </>
              ) : (
                <>
                  <LoginModal
                    open={showLoginModal}
                    onOpenChange={setShowLoginModal}
                    onSwitchToRegister={handleSwitchToRegister}
                  >
                    <Button variant="ghost" className="font-manrope text-white hover:bg-white/10">
                      Войти
                    </Button>
                  </LoginModal>
                  <RegistrationModal
                    open={showRegistrationModal}
                    onOpenChange={setShowRegistrationModal}
                    onSwitchToLogin={handleSwitchToLogin}
                  >
                    <Button className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600">
                      Регистрация
                    </Button>
                  </RegistrationModal>
                </>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="relative z-20 pt-24 pb-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Page Header */}
          <div className="text-center mb-12">
            <h1 className="font-geist font-bold text-4xl lg:text-5xl text-white mb-4">История обработки</h1>
            <p className="font-manrope text-xl text-white/70 max-w-2xl mx-auto">
              Все ваши обработанные файлы в одном месте. Скачивайте повторно или просматривайте детали обработки.
            </p>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
            <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
              <CardContent className="p-6 text-center">
                <div className="w-12 h-12 bg-gradient-to-r from-purple-500 to-blue-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <FileImage className="w-6 h-6 text-white" />
                </div>
                <div className="text-2xl font-bold mb-1">{stats.total_files}</div>
                <div className="text-white/70 text-sm">Обработано файлов</div>
              </CardContent>
            </Card>
            <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
              <CardContent className="p-6 text-center">
                <div className="w-12 h-12 bg-gradient-to-r from-green-500 to-teal-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <TrendingUp className="w-6 h-6 text-white" />
                </div>
                <div className="text-2xl font-bold mb-1">{stats.processed_this_week}</div>
                <div className="text-white/70 text-sm">Обработано за неделю</div>
              </CardContent>
            </Card>
            <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
              <CardContent className="p-6 text-center">
                <div className="w-12 h-12 bg-gradient-to-r from-orange-500 to-red-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Clock className="w-6 h-6 text-white" />
                </div>
                <div className="text-2xl font-bold mb-1">{(stats.total_files * 0.2).toFixed(1)}</div>
                <div className="text-white/70 text-sm">Часа сэкономлено</div>
              </CardContent>
            </Card>
          </div>

          {/* Files List */}
          {isAuthenticated ? (
            loading ? (
              <div className="text-center py-16">
                <div className="w-24 h-24 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-6 animate-pulse">
                  <FileImage className="w-12 h-12 text-white/40" />
                </div>
                <h3 className="font-geist font-semibold text-xl text-white mb-2">Загрузка истории...</h3>
              </div>
            ) : files.length === 0 ? (
              <div className="text-center py-16">
                <div className="w-24 h-24 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-6">
                  <FileImage className="w-12 h-12 text-white/40" />
                </div>
                <h3 className="font-geist font-semibold text-xl text-white mb-2">Пока нет обработанных файлов</h3>
                <p className="text-white/70 mb-6">Загрузите и обработайте свой первый файл, чтобы увидеть его здесь</p>
                <Button
                  asChild
                  className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                >
                  <Link href="/process">Обработать файл</Link>
                </Button>
              </div>
            ) : (
              <div className="space-y-4">
                {files.map((file) => (
                  <Card
                    key={file.id}
                    className="bg-white/5 backdrop-blur-lg border-white/10 text-white hover:bg-white/10 transition-all duration-300"
                  >
                    <CardContent className="p-6">
                      <div className="flex items-center gap-6">
                        {/* Thumbnail */}
                        <div className="flex-shrink-0">
                          <div className="w-20 h-20 bg-white/10 rounded-lg overflow-hidden flex items-center justify-center">
                            {file.thumbnail ? (
                              <img
                                src={file.thumbnail}
                                alt={file.original_name}
                                className="w-full h-full object-cover"
                              />
                            ) : getFileType(file.mime_type) === "image" ? (
                              <FileImage className="w-8 h-8 text-white/60" />
                            ) : (
                              <FileVideo className="w-8 h-8 text-white/60" />
                            )}
                          </div>
                        </div>

                        {/* File Info */}
                        <div className="flex-grow">
                          <div className="flex items-center gap-3 mb-2">
                            <h3 className="font-manrope font-semibold text-lg text-white">{file.original_name}</h3>
                            <Badge variant="secondary" className="bg-white/10 text-white border-white/20">
                              {getFileType(file.mime_type) === "image" ? "Фото" : "Видео"}
                            </Badge>
                          </div>

                          <div className="grid grid-cols-2 md:grid-cols-2 gap-2 text-sm text-white/70">
                            <div className="flex items-center gap-2">
                              <Calendar className="w-4 h-4" />
                              <span>{formatDate(file.processed_at || file.uploaded_at)}</span>
                            </div>
                            <div>
                              <span className="font-medium">Размер:</span>{" "}
                              {formatFileSize(file.processed_size || file.file_size)}
                            </div>
                            {/*
                            <div>
                              <span className="font-medium">Эффект:</span>{" "}
                              {blurTypeMap[file.blur_type] || "Неизвестно"}
                            </div>
                            <div>
                              <span className="font-medium">Качество:</span>{" "}
                              {qualityMap[file.quality] || "Оригинальное"}
                            </div>
                            */}
                          </div>

                          <div className="mt-3">
                            <div className="flex flex-wrap gap-2">
                              <span className="text-white/70 text-sm">Статус:</span>
                              <Badge variant="outline" className="border-cyan-400/30 text-cyan-400 text-xs">
                                {statusMap[file.status] || file.status}
                              </Badge>
                              {file.objects_detected?.length > 0 && (
                                <>
                                  <span className="text-white/70 text-sm">Размыто:</span>
                                  {file.objects_detected.map((object: string, index: number) => (
                                    <Badge
                                      key={index}
                                      variant="outline"
                                      className="border-cyan-400/30 text-cyan-400 text-xs"
                                    >
                                      {object}
                                    </Badge>
                                  ))}
                                </>
                              )}
                              {file.status === "failed" && file.error_message && (
                                <span className="text-red-400 text-xs">{file.error_message}</span>
                              )}
                            </div>
                          </div>
                        </div>

                        {/* Actions */}
                        <div className="flex-shrink-0 flex gap-2">
                          {file.thumbnail && (
                            <Button
                              variant="outline"
                              size="sm"
                              className="border-white/20 text-white hover:bg-white/10 bg-transparent"
                              onClick={() => window.open(file.thumbnail, "_blank")}
                            >
                              <Eye className="w-4 h-4 mr-2" />
                              Просмотр
                            </Button>
                          )}
                          {file.status === "completed" && (
                            <>
                              <Button
                                size="sm"
                                className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                                onClick={() => handleDownload(file.id, "original")}
                              >
                                <Download className="w-4 h-4 mr-2" />
                                Оригинал
                              </Button>
                              <Button
                                size="sm"
                                className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                                onClick={() => handleDownload(file.id, "processed")}
                              >
                                <Download className="w-4 h-4 mr-2" />
                                Обработанный
                              </Button>
                            </>
                          )}
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDelete(file.id)}
                          >
                            <Trash className="w-4 h-4 mr-2" />
                            Удалить
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )
          ) : (
            <div className="text-center py-16">
              <div className="w-24 h-24 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-6">
                <FileImage className="w-12 h-12 text-white/40" />
              </div>
              <h3 className="font-geist font-semibold text-xl text-white mb-2">Войдите, чтобы увидеть историю</h3>
              <p className="text-white/70 mb-6">Авторизуйтесь, чтобы просмотреть ваши обработанные файлы</p>
              <div className="flex justify-center gap-4">
                <LoginModal
                  open={showLoginModal}
                  onOpenChange={setShowLoginModal}
                  onSwitchToRegister={handleSwitchToRegister}
                >
                  <Button variant="ghost" className="font-manrope text-white hover:bg-white/10">
                    Войти
                  </Button>
                </LoginModal>
                <RegistrationModal
                  open={showRegistrationModal}
                  onOpenChange={setShowRegistrationModal}
                  onSwitchToLogin={handleSwitchToLogin}
                >
                  <Button className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600">
                    Регистрация
                  </Button>
                </RegistrationModal>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}