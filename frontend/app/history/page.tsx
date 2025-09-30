"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Camera, Download, TrendingUp, Calendar, CircleX, FileImage, FileVideo, ArrowLeft, Trash, Menu, User, Settings, LogOut, ChevronDown, ChevronUp } from "lucide-react"
import Link from "next/link"
import { useAuth } from "@/components/AuthContext"
import { LoginModal } from "@/components/login-modal"
import { RegistrationModal } from "@/components/registration-modal"
import CameraIcon from "@/components/ui/camera-icon";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"

const API_LINK = process.env.NEXT_PUBLIC_API_LINK || "http://localhost:8080";

export default function HistoryPage() {
  const { token, isAuthenticated, user, logout } = useAuth()
  const [files, setFiles] = useState<any[]>([])
  const [stats, setStats] = useState<{ total_files: number; total_processed: number; total_failed: number; total_size: number }>({
    total_files: 0,
    total_processed: 0,
    total_failed: 0,
    total_size: 0,
  })
  const [loading, setLoading] = useState(true)
  const [showLoginModal, setShowLoginModal] = useState(false)
  const [showRegistrationModal, setShowRegistrationModal] = useState(false)
  const [open, setOpen] = useState(false)
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [selectedMedia, setSelectedMedia] = useState<{ url: string; type: "image" | "video" } | null>(null)
  const [thumbnails, setThumbnails] = useState<{ [fileId: string]: string }>({})

  const truncateFileName = (name: string, maxLength: number = 20) => {
    if (name.length <= maxLength) return name;
    const extension = name.split(".").pop() || "";
    const nameWithoutExt = name.substring(0, name.lastIndexOf("."));
    const truncatedName = nameWithoutExt.substring(0, maxLength - 3 - extension.length);
    return `${truncatedName}...${extension}`;
  };

  const openFullScreen = (url: string, type: "image" | "video") => {
    setSelectedMedia({ url, type });
  };

  const closeFullScreen = () => {
    setSelectedMedia(null);
  };

  const generateVideoThumbnail = (videoUrl: string, fileId: string) => {
    return new Promise<string>((resolve) => {
      const video = document.createElement("video");
      video.src = videoUrl;
      video.crossOrigin = "anonymous";
      video.muted = true;
      video.preload = "metadata";

      video.onloadedmetadata = () => {
        video.currentTime = 1;
      };

      video.onseeked = () => {
        const canvas = document.createElement("canvas");
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        const ctx = canvas.getContext("2d");
        if (ctx) {
          try {
            ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
            const dataUrl = canvas.toDataURL("image/jpeg");
            setThumbnails((prev) => ({ ...prev, [fileId]: dataUrl }));
            resolve(dataUrl);
          } catch (error) {
            console.error(`Ошибка генерации превью для ${fileId}:`, error);
            resolve("");
          }
        } else {
          resolve("");
        }
        video.remove();
      };

      video.onerror = () => {
        console.error(`Ошибка загрузки видео ${fileId}:`, video.error);
        resolve("");
        video.remove();
      };
    });
  };

  const fetchFiles = async () => {
    if (!isAuthenticated || !token) {
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const res = await fetch(`${API_LINK}/api/files`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) {
        if (res.status === 401 || res.status === 403) {
          alert("Сессия истекла, пожалуйста, войдите снова");
          logout();
          return;
        }
        throw new Error(`Ошибка получения файлов: ${res.status}`);
      }
      const { data } = await res.json();
      const filesWithUrls = await Promise.all(
        data.map(async (file: any) => {
          if (file.status === "completed") {
            const processedUrl = `${API_LINK}/api/files/${file.id}?type=processed`;
            if (getFileType(file.mime_type) === "video") {
              const thumbnail = await generateVideoThumbnail(processedUrl, file.id).catch((error) => {
                console.error(`Ошибка генерации превью для ${file.id}:`, error);
                return "";
              });
              return { ...file, processed_url: processedUrl, thumbnail };
            }
            return { ...file, processed_url: processedUrl };
          }
          return file;
        })
      );
      setFiles(filesWithUrls || []);
    } catch (error) {
      console.error("Ошибка загрузки файлов:", error);
      alert("Не удалось загрузить историю файлов");
    } finally {
      setLoading(false);
    }
  };

  const fetchStats = async () => {
    if (!isAuthenticated || !token) return;
    try {
      const res = await fetch(`${API_LINK}/api/user/stats`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) {
        if (res.status === 401 || res.status === 403) {
          alert("Сессия истекла, пожалуйста, войдите снова");
          logout();
          return;
        }
        throw new Error(`Ошибка получения статистики: ${res.status}`);
      }
      const { data } = await res.json();
      setStats({
        total_files: data.total_files || 0,
        total_processed: data.total_processed || 0,
        total_failed: data.total_failed || 0,
        total_size: data.total_size || 0,
      });
    } catch (error) {
      console.error("Ошибка загрузки статистики:", error);
      setStats({ total_files: 0, total_processed: 0, total_failed: 0, total_size: 0 });
    }
  };

  const handleDownload = async (fileId: string, type: "original" | "processed") => {
    try {
      const res = await fetch(`${API_LINK}/api/files/${fileId}?type=${type}`, {
        method: "GET",
        headers: isAuthenticated && token ? { Authorization: `Bearer ${token}` } : undefined,
      });

      if (!res.ok) {
        const error = await res.json().catch(() => ({}));
        if ((res.status === 401 || res.status === 403) && isAuthenticated) {
          alert("Сессия истекла, пожалуйста, войдите снова");
          logout();
          return;
        }
        throw new Error(error.message || `Ошибка скачивания: ${res.status}`);
      }

      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;

      const file = files.find((f) => f.id === fileId);
      const extension = file?.original_name
        ? file.original_name.split(".").pop() || "file"
        : "file";
      a.download = `${type}-${fileId}.${extension}`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (error: any) {
      console.error(`Ошибка скачивания (${type}):`, error);
      alert(`Ошибка: ${error.message || "Не удалось скачать файл"}`);
    }
  };

  const handleDelete = async (fileId: string) => {
    if (!isAuthenticated || !token) return;
    try {
      const res = await fetch(`${API_LINK}/api/files/${fileId}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) {
        if (res.status === 401 || res.status === 403) {
          alert("Сессия истекла, пожалуйста, войдите снова");
          logout();
          return;
        }
        throw new Error(`Ошибка удаления: ${res.status}`);
      }
      setFiles(files.filter((f) => f.id !== fileId));
      setThumbnails((prev) => {
        const newThumbnails = { ...prev };
        delete newThumbnails[fileId];
        return newThumbnails;
      });
    } catch (error) {
      console.error("Ошибка удаления:", error);
      alert("Не удалось удалить файл");
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getFileType = (mimeType: string) => {
    return mimeType.startsWith("image/") ? "image" : mimeType.startsWith("video/") ? "video" : "file";
  };

  const formatFileSize = (size: number) => {
    if (size < 1024) return `${size} Б`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} КБ`;
    return `${(size / (1024 * 1024)).toFixed(1)} МБ`;
  };

  const formatTotalSizeMB = (size: number) => {
    return (size / (1024 * 1024)).toFixed(1);
  };

  const statusMap: { [key: string]: string } = {
    uploaded: "Загружен",
    processing: "Обрабатывается",
    completed: "Завершено",
    failed: "Ошибка",
  };

  useEffect(() => {
    if (isAuthenticated && token) {
      fetchFiles();
      fetchStats();
    }

    const cursor = document.querySelector(".cursor");
    const cursorFollower = document.querySelector(".cursor-follower");
    let cursorX = 0;
    let cursorY = 0;

    const handleMouseMove = (e: MouseEvent) => {
      cursorX = e.clientX;
      cursorY = e.clientY;

      if (cursor) {
        cursor.style.transform = `translate(${cursorX}px, ${cursorY}px)`;
      }

      if (cursorFollower) {
        cursorFollower.style.transform = `translate(${cursorX}px, ${cursorY}px)`;
      }
    };

    const handleMouseDown = () => {
      cursor?.classList.add("active");
    };

    const handleMouseUp = () => {
      cursor?.classList.remove("active");
    };

    const createParticles = () => {
      const particlesContainer = document.getElementById("particles");
      if (!particlesContainer) return;

      const particleCount = 30;

      for (let i = 0; i < particleCount; i++) {
        const particle = document.createElement("div");
        particle.className = "particle";
        particle.style.left = Math.random() * 100 + "%";
        particle.style.animationDelay = Math.random() * 20 + "s";
        particle.style.animationDuration = Math.random() * 10 + 10 + "s";
        particle.style.width = Math.random() * 4 + 1 + "px";
        particle.style.height = particle.style.width;
        particle.style.animation = `particle-up ${particle.style.animationDuration} linear infinite`;
        particlesContainer.appendChild(particle);
      }
    };

    const handleScroll = () => {
      const header = document.getElementById("header");
      if (header) {
        if (window.scrollY > 50) {
          header.classList.add("scrolled");
        } else {
          header.classList.remove("scrolled");
        }
      }
    };

    const setupToggleButtons = () => {
      document.querySelectorAll(".toggle-group").forEach((group) => {
        const buttons = group.querySelectorAll(".toggle-btn");
        buttons.forEach((btn) => {
          btn.addEventListener("click", () => {
            buttons.forEach((b) => b.classList.remove("active"));
            btn.classList.add("active");
          });
        });
      });
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mousedown", handleMouseDown);
    document.addEventListener("mouseup", handleMouseUp);
    window.addEventListener("scroll", handleScroll);
    
    createParticles();
    setupToggleButtons();

    if (cursor) cursor.style.opacity = "1";
    if (cursorFollower) cursorFollower.style.opacity = "1";

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mousedown", handleMouseDown);
      document.removeEventListener("mouseup", handleMouseUp);
      window.removeEventListener("scroll", handleScroll);
    };
  }, [isAuthenticated, token]);

  const handleSwitchToRegister = () => {
    setShowLoginModal(false);
    setTimeout(() => setShowRegistrationModal(true), 100);
  };

  const handleSwitchToLogin = () => {
    setShowRegistrationModal(false);
    setTimeout(() => setShowLoginModal(true), 100);
  };

  return (
    <div className="min-h-screen bg-black relative">
      <div className="cursor"></div>
      <div className="cursor-follower"></div>

      <div className="bg-animation absolute top-0 left-0 w-full h-full z-0"></div>
      <div className="particles absolute top-0 left-0 w-full overflow-hidden h-full z-10" id="particles"></div>

      <header className="fixed top-0 w-full z-50 backdrop-blur-lg bg-black/10 border-b border-white/10 transition-all duration-300">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-4">
              <Link href="/" className="flex items-center gap-2 text-white/80 hover:text-white transition-colors">
                <ArrowLeft className="w-5 h-5" />
              </Link>
              <Link href="/" className="logo">
                <div className="logo-icon">
                  <CameraIcon />
                </div>
                <span className="logo-text">Obscura</span>
              </Link>
            </div>

            <DropdownMenu open={isMenuOpen} onOpenChange={setIsMenuOpen}>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="outline"
                  className="flex items-center gap-2 font-manrope text-white bg-white/10 hover:bg-white/20 border-0 md:hidden"
                >
                  <Menu className="h-5 w-5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56 z-2000">
                <div className="flex flex-col gap-2 p-2">
                  <DropdownMenuItem className="hover:bg-transparent focus:bg-transparent text-black hover:text-black focus:text-black">
                    <Link href="/" onClick={() => setIsMenuOpen(false)}>
                      Домой
                    </Link>
                  </DropdownMenuItem>
                  {isAuthenticated && (
                    <DropdownMenuItem className="hover:bg-transparent focus:bg-transparent text-black hover:text-black focus:text-black">
                      <Link href="/process" onClick={() => setIsMenuOpen(false)}>
                        Обработать
                      </Link>
                    </DropdownMenuItem>
                  )}
                  <DropdownMenuSeparator />
                  {isAuthenticated ? (
                    <>
                      <DropdownMenuLabel>
                        <div className="flex flex-col">
                          <span className="font-semibold">{user?.name || "Пользователь"}</span>
                          <span className="text-sm text-[#8c939f]">{user?.email}</span>
                        </div>
                      </DropdownMenuLabel>
                      <DropdownMenuItem className="hover:bg-transparent focus:bg-transparent text-black hover:text-black focus:text-black">
                        <Link href="/dashboard" className="flex items-center" onClick={() => setIsMenuOpen(false)}>
                          <Settings className="mr-2 h-4 w-4 text-current" />
                          Настройки
                        </Link>
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => {
                          logout();
                          setIsMenuOpen(false);
                        }}
                        className="text-red-500 hover:bg-transparent focus:bg-transparent focus:text-red-500"
                      >
                        <LogOut className="mr-2 h-4 w-4 text-current" />
                        Выйти
                      </DropdownMenuItem>
                    </>
                  ) : (
                    <>
                      <DropdownMenuItem className="hover:bg-transparent focus:bg-transparent text-black hover:text-black focus:text-black">
                        <LoginModal
                          open={showLoginModal}
                          onOpenChange={(open) => {
                            setShowLoginModal(open);
                            setIsMenuOpen(false);
                          }}
                          onSwitchToRegister={() => {
                            handleSwitchToRegister();
                            setIsMenuOpen(false);
                          }}
                        >
                          <button className="w-full text-left">Войти</button>
                        </LoginModal>
                      </DropdownMenuItem>
                      <DropdownMenuItem className="hover:bg-transparent focus:bg-transparent text-black hover:text-black focus:text-black">
                        <RegistrationModal
                          open={showRegistrationModal}
                          onOpenChange={(open) => {
                            setShowRegistrationModal(open);
                            setIsMenuOpen(false);
                          }}
                          onSwitchToLogin={() => {
                            handleSwitchToLogin();
                            setIsMenuOpen(false);
                          }}
                        >
                          <button className="w-full text-left">Регистрация</button>
                        </RegistrationModal>
                      </DropdownMenuItem>
                    </>
                  )}
                </div>
              </DropdownMenuContent>
            </DropdownMenu>

            <div className="hidden md:flex items-center gap-3">
              {isAuthenticated && (
                <Link href="/process" className="text-white/80 hover:text-white transition-colors">
                  <Button className="font-manrope bg-gradient-to-r from-cyan-500/20 to-blue-500/20 border border-cyan-400/30 text-white hover:from-cyan-500/30 hover:to-blue-500/30 hover:border-cyan-400/50 transition-all duration-300 transform hover:scale-105 hover:shadow-lg hover:shadow-cyan-500/25 backdrop-blur-sm flex items-center gap-2">
                    <Camera className="w-4 h-4" />
                    Обработать
                  </Button>
                </Link>
              )}
              {isAuthenticated ? (
                <DropdownMenu open={open} onOpenChange={setOpen}>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      className="flex items-center gap-2 font-manrope text-white bg-white/10 hover:bg-white/20 border-0"
                    >
                      <User className="h-4 w-4" />
                      {user?.name || user?.email || "Пользователь"}
                      {open ? <ChevronUp className="h-4 w-4 ml-1" /> : <ChevronDown className="h-4 w-4 ml-1" />}
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56 z-2000">
                    <DropdownMenuLabel>
                      <div className="flex flex-col">
                        <span className="font-semibold">{user?.name || "Пользователь"}</span>
                        <span className="text-sm text-[#8c939f]">{user?.email}</span>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <div className="flex flex-col gap-4 mt-4 pl-3">
                      <DropdownMenuItem className="hover:bg-transparent focus:bg-transparent text-black hover:text-black focus:text-black" style={{ fontSize: "17px" }}>
                        <Link href="/dashboard" className="flex items-center">
                          <Settings className="mr-2 h-4 w-4 text-current" />
                          Настройки
                        </Link>
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={logout}
                        className="text-red-500 hover:bg-transparent focus:bg-transparent focus:text-red-500"
                        style={{ fontSize: "17px" }}
                      >
                        <LogOut className="mr-2 h-4 w-4 text-current" />
                        Выйти
                      </DropdownMenuItem>
                    </div>
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : (
                <>
                  <LoginModal open={showLoginModal} onOpenChange={setShowLoginModal} onSwitchToRegister={handleSwitchToRegister}>
                    <Button variant="ghost" className="font-manrope text-white hover:bg-white/10">
                      Войти
                    </Button>
                  </LoginModal>
                  <RegistrationModal open={showRegistrationModal} onOpenChange={setShowRegistrationModal} onSwitchToLogin={handleSwitchToLogin}>
                    <Button className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 shadow-lg hover:shadow-purple-500/25 transition-all duration-300">
                      Регистрация
                    </Button>
                  </RegistrationModal>
                </>
              )}
            </div>
          </div>
        </div>
      </header>

      <main className="relative z-20 pt-24 pb-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h1 className="font-geist font-bold text-4xl lg:text-5xl text-white mb-4">История обработки</h1>
            <p className="font-manrope text-xl text-white/70 max-w-2xl mx-auto">
              Все ваши обработанные файлы в одном месте. Статистика сохраняется даже после удаления файлов.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-12">
            <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
              <CardContent className="p-6 text-center">
                <div className="w-12 h-12 bg-gradient-to-r from-purple-500 to-blue-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <FileImage className="w-6 h-6 text-white" />
                </div>
                <div className="text-2xl font-bold mb-1">{stats.total_files}</div>
                <div className="text-white/70 text-sm">Всего файлов</div>
              </CardContent>
            </Card>
            <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
              <CardContent className="p-6 text-center">
                <div className="w-12 h-12 bg-gradient-to-r from-green-500 to-teal-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <TrendingUp className="w-6 h-6 text-white" />
                </div>
                <div className="text-2xl font-bold mb-1">{stats.total_processed}</div>
                <div className="text-white/70 text-sm">Успешно обработано</div>
              </CardContent>
            </Card>
            <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
              <CardContent className="p-6 text-center">
                <div className="w-12 h-12 bg-gradient-to-r from-orange-500 to-red-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <CircleX className="w-6 h-6 text-white" />
                </div>
                <div className="text-2xl font-bold mb-1">{stats.total_failed}</div>
                <div className="text-white/70 text-sm">Неудачных обработок</div>
              </CardContent>
            </Card>
            <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
              <CardContent className="p-6 text-center">
                <div className="w-12 h-12 bg-gradient-to-r from-cyan-500 to-blue-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Download className="w-6 h-6 text-white" />
                </div>
                <div className="text-2xl font-bold mb-1">{formatTotalSizeMB(stats.total_size)}</div>
                <div className="text-white/70 text-sm">МБ обработано</div>
              </CardContent>
            </Card>
          </div>

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
              <div className="space-y-4 max-w-full overflow-x-hidden">
                {files.map((file) => (
                  <Card
                    key={file.id}
                    className="bg-white/5 backdrop-blur-lg border-white/10 text-white hover:bg-white/10 transition-all duration-300 w-full max-w-full"
                  >
                    <CardContent className="p-6">
                      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:gap-6">
                        <div className="flex-shrink-0 w-16 h-16 sm:w-20 sm:h-20">
                          <div className="w-full h-full bg-white/10 rounded-lg overflow-hidden flex items-center justify-center">
                            {file.status === "completed" && file.processed_url ? (
                              getFileType(file.mime_type) === "image" ? (
                                <img
                                  src={file.processed_url}
                                  alt={file.original_name}
                                  className="w-full h-full object-cover cursor-pointer"
                                  onClick={() => openFullScreen(file.processed_url, "image")}
                                />
                              ) : thumbnails[file.id] ? (
                                <img
                                  src={thumbnails[file.id]}
                                  alt={`${file.original_name} thumbnail`}
                                  className="w-full h-full object-cover cursor-pointer"
                                  onClick={() => openFullScreen(file.processed_url, "video")}
                                />
                              ) : (
                                <FileVideo className="w-6 h-6 sm:w-8 sm:h-8 text-white/60" />
                              )
                            ) : getFileType(file.mime_type) === "image" ? (
                              <FileImage className="w-6 h-6 sm:w-8 sm:h-8 text-white/60" />
                            ) : (
                              <FileVideo className="w-6 h-6 sm:w-8 sm:h-8 text-white/60" />
                            )}
                          </div>
                        </div>

                        <div className="flex-grow">
                          <div className="flex items-center gap-3 mb-2">
                            <h3 className="font-manrope font-semibold text-lg text-white">{truncateFileName(file.original_name)}</h3>
                            <Badge variant="secondary" className="bg-white/10 text-white border-white/20">
                              {getFileType(file.mime_type) === "image" ? "Фото" : "Видео"}
                            </Badge>
                          </div>

                          <div className="flex flex-col sm:grid sm:grid-cols-2 gap-2 text-sm text-white/70">
                            <div className="flex items-center gap-2">
                              <Calendar className="w-4 h-4" />
                              <span>{formatDate(file.processed_at || file.uploaded_at)}</span>
                            </div>
                            <div>
                              <span className="font-medium">Размер:</span>{" "}
                              {formatFileSize(file.processed_size || file.file_size)}
                            </div>
                          </div>

                          <div className="mt-3">
                            <div className="flex flex-wrap gap-2">
                              <span className="text-white/70 text-sm">Статус:</span>
                              <Badge variant="outline" className="border-cyan-400/30 text-cyan-400 text-xs">
                                {statusMap[file.status] || file.status}
                              </Badge>
                              {file.status === "failed" && file.error_message && (
                                <span className="text-red-400 text-xs">{file.error_message}</span>
                              )}
                            </div>
                          </div>
                        </div>

                        <div className="flex-shrink-0 flex flex-wrap gap-2 mt-4 sm:mt-0">
                          {file.status === "completed" && (
                            <>
                              <Button
                                size="sm"
                                className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 text-xs sm:text-sm px-2 sm:px-4"
                                onClick={() => handleDownload(file.id, "original")}
                              >
                                <Download className="w-4 h-4 mr-1 sm:mr-2" />
                                Оригинал
                              </Button>
                              <Button
                                size="sm"
                                className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 text-xs sm:text-sm px-2 sm:px-4"
                                onClick={() => handleDownload(file.id, "processed")}
                              >
                                <Download className="w-4 h-4 mr-1 sm:mr-2" />
                                Обработанный
                              </Button>
                            </>
                          )}
                          <Button
                            size="sm"
                            variant="destructive"
                            className="text-xs sm:text-sm px-2 sm:px-4"
                            onClick={() => handleDelete(file.id)}
                          >
                            <Trash className="w-4 h-4 mr-1 sm:mr-2" />
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

          {selectedMedia && (
            <div
              className="fixed inset-0 bg-black/90 z-50 flex items-center justify-center"
              onClick={closeFullScreen}
            >
              {selectedMedia.type === "image" ? (
                <img
                  src={selectedMedia.url}
                  alt="Full screen preview"
                  className="max-w-[90%] max-h-[90%] object-contain"
                />
              ) : (
                <video
                  src={selectedMedia.url}
                  controls
                  autoPlay
                  className="max-w-[90%] max-h-[90%] object-contain"
                />
              )}
              <Button
                variant="ghost"
                className="absolute top-4 right-4 text-white hover:bg-white/10"
                onClick={closeFullScreen}
              >
                <CircleX className="w-6 h-6" />
              </Button>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}