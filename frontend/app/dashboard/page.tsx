"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Progress } from "@/components/ui/progress";
import { Input } from "@/components/ui/input";
import {
  Camera,
  Settings,
  Crown,
  FileImage,
  Video,
  Clock,
  Download,
  Edit3,
  Calendar,
  CreditCard,
  Shield,
  Bell,
  User,
  ChevronRight,
  CircleX,
} from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/components/AuthContext";
import { LoginModal } from "@/components/login-modal";
import { RegistrationModal } from "@/components/registration-modal";

const API_LINK = process.env.NEXT_PUBLIC_API_LINK || "http://localhost:8080";

export default function DashboardPage() {
  const { token, isAuthenticated, user, logout } = useAuth();
  const [stats, setStats] = useState({
    total_files: 0,
    total_processed: 0,
    total_failed: 0,
    total_size: 0,
    created_at: "",
    updated_at: "",
    last_stats_update: "",
  });
  const [loading, setLoading] = useState(true);
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [showRegistrationModal, setShowRegistrationModal] = useState(false);
  const [editing, setEditing] = useState(false);
  const [profileData, setProfileData] = useState({
    name: "",
    email: "",
    password: "",
  });
  const [initialProfileData, setInitialProfileData] = useState({
    name: "",
    email: "",
    password: "",
  });
  const [errors, setErrors] = useState({
    name: "",
    email: "",
    password: "",
  });

  // Загрузка статистики пользователя
  const fetchStats = async () => {
    if (!isAuthenticated || !token) {
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const res = await fetch(`${API_LINK}/api/user/stats`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) {
        if (res.status === 401) {
          alert("Сессия истекла, пожалуйста, войдите снова");
          logout();
        }
        throw new Error(`Ошибка получения статистики: ${res.status}`);
      }
      const { data } = await res.json();
      setStats({
        total_files: data.total_files || 0,
        total_processed: data.total_processed || 0,
        total_failed: data.total_failed || 0,
        total_size: data.total_size || 0,
        created_at: data.created_at || "",
        updated_at: data.updated_at || "",
        last_stats_update: data.last_stats_update || "",
      });
      const profile = {
        name: data.name || "",
        email: data.email || "",
        password: "",
      };
      setProfileData(profile);
      setInitialProfileData(profile);
    } catch (error) {
      console.error("Ошибка загрузки статистики:", error);
      setStats({
        total_files: 0,
        total_processed: 0,
        total_failed: 0,
        total_size: 0,
        created_at: "",
        updated_at: "",
        last_stats_update: "",
      });
      setProfileData({ name: "", email: "", password: "" });
      setInitialProfileData({ name: "", email: "", password: "" });
    } finally {
      setLoading(false);
    }
  };

  // Валидация полей
  const validateFields = () => {
    const newErrors = { name: "", email: "", password: "" };
    let isValid = true;

    if (profileData.name && profileData.name.length < 2) {
      newErrors.name = "Имя должно содержать минимум 2 символа";
      isValid = false;
    }

    if (profileData.email) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(profileData.email)) {
        newErrors.email = "Введите корректный email";
        isValid = false;
      }
    }

    if (profileData.password && profileData.password.length < 6) {
      newErrors.password = "Пароль должен содержать минимум 6 символов";
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  // Обновление профиля
  const handleUpdateProfile = async () => {
    if (!isAuthenticated || !token) return;

    if (!validateFields()) {
      return;
    }

    try {
      const body = {};
      if (profileData.name) body.name = profileData.name;
      if (profileData.email) body.email = profileData.email;
      if (profileData.password) body.password = profileData.password;

      if (Object.keys(body).length === 0) {
        setErrors({ ...errors, email: "Введите хотя бы одно поле для обновления" });
        return;
      }

      const res = await fetch(`${API_LINK}/api/user/profile/update`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (!res.ok) {
        if (res.status === 401) {
          alert("Сессия истекла, пожалуйста, войдите снова");
          logout();
        }
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error || `Ошибка обновления профиля: ${res.status}`);
      }

      const { data } = await res.json();
      const updatedProfile = {
        name: data.name || "",
        email: data.email || "",
        password: "",
      };
      setProfileData(updatedProfile);
      setInitialProfileData(updatedProfile);
      setStats((prev) => ({
        ...prev,
        created_at: data.created_at || prev.created_at,
        updated_at: data.updated_at || prev.updated_at,
        last_stats_update: data.last_stats_update || prev.last_stats_update,
      }));
      setEditing(false);
      setErrors({ name: "", email: "", password: "" });
      alert("Профиль успешно обновлен");
    } catch (error) {
      console.error("Ошибка обновления профиля:", error);
      setErrors({ ...errors, email: error.message || "Не удалось обновить профиль" });
    }
  };

  // Форматирование размера файла
  const formatTotalSizeMB = (size) => {
    return (size / (1024 * 1024)).toFixed(1);
  };

  // Форматирование даты
  const formatDate = (dateString) => {
    if (!dateString) return "Неизвестно";
    const date = new Date(dateString);
    return date.toLocaleDateString("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
    });
  };

  // Обработка переключения на регистрацию
  const handleSwitchToRegister = () => {
    setShowLoginModal(false);
    setTimeout(() => setShowRegistrationModal(true), 100);
  };

  // Обработка переключения на логин
  const handleSwitchToLogin = () => {
    setShowRegistrationModal(false);
    setTimeout(() => setShowLoginModal(true), 100);
  };

  // Инициализация эффектов
  useEffect(() => {
    if (isAuthenticated && token) {
      fetchStats();
    }

    // Логика курсора
    const cursor = document.querySelector(".cursor");
    const cursorFollower = document.querySelector(".cursor-follower");
    let cursorX = 0;
    let cursorY = 0;

    const handleMouseMove = (e) => {
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

    // Создание частиц
    const createParticles = () => {
      const particlesContainer = document.getElementById("particles");
      if (!particlesContainer) return;

      particlesContainer.innerHTML = "";

      const particleCount = 30;

      for (let i = 0; i < particleCount; i++) {
        const particle = document.createElement("div");
        particle.className = "particle";
        particle.style.width = Math.random() * 4 + 1 + "px";
        particle.style.height = particle.style.width;
        particle.style.left = Math.random() * 100 + "%";
        particle.style.animationDuration = Math.random() * 20 + 10 + "s";
        particle.style.animationDelay = Math.random() * 20 + "s";
        particle.style.animation = `particle-up ${particle.style.animationDuration} linear infinite`;
        particlesContainer.appendChild(particle);
      }
    };

    // Логика прокрутки заголовка
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

    // Логика переключателей
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

    // Инициализация видимости .feature-card
    const cards = document.querySelectorAll(".feature-card");
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add("visible");
          }
        });
      },
      { threshold: 0.1 }
    );

    cards.forEach((card) => {
      observer.observe(card);
    });

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mousedown", handleMouseDown);
    document.addEventListener("mouseup", handleMouseUp);
    window.addEventListener("scroll", handleScroll);
    createParticles();
    setupToggleButtons();

    if (cursor) cursor.style.opacity = "1";
    if (cursorFollower) cursorFollower.style.opacity = "1";

    // Очистка
    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mousedown", handleMouseDown);
      document.removeEventListener("mouseup", handleMouseUp);
      window.removeEventListener("scroll", handleScroll);
      observer.disconnect();
    };
  }, [isAuthenticated, token]);

  return (
    <div className="min-h-screen bg-black relative">
      <div className="cursor"></div>
      <div className="cursor-follower"></div>

      {/* Background Animation */}
      <div className="bg-animation absolute top-0 left-0 w-full h-full z-0"></div>

      {/* Header */}
      <header className="fixed top-0 w-full z-50 backdrop-blur-lg bg-black/10 border-b border-white/10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link href="/" className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-br from-white to-purple-400 rounded-lg flex items-center justify-center">
                <Camera className="w-5 h-5 text-black" />
              </div>
              <span className="font-geist font-bold text-xl bg-gradient-to-r from-white to-purple-400 bg-clip-text text-transparent">
                Obscura
              </span>
            </Link>
            <nav className="hidden md:flex items-center gap-6">
              <Link href="/" className="text-white/80 hover:text-white transition-colors">
                Главная
              </Link>
              <Link href="/process" className="text-white/80 hover:text-white transition-colors">
                Обработка
              </Link>
              <Link href="/history" className="text-white/80 hover:text-white transition-colors">
                История
              </Link>
            </nav>
            <div className="flex items-center gap-3">
              {isAuthenticated ? (
                <Button
                  variant="outline"
                  className="bg-white/5 border-white/20 text-white hover:bg-white/10"
                  onClick={logout}
                >
                  <User className="w-4 h-4 mr-2" />
                  Выйти
                </Button>
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
      <div className="relative z-20 pt-24 pb-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {isAuthenticated ? (
            <>
              {/* Welcome Section */}
              <div className="mb-8">
                <div className="flex items-center gap-4 mb-6">
                  <Avatar className="w-16 h-16 border-2 border-white/20">
                    <AvatarImage src="/placeholder-user.png" />
                    <AvatarFallback className="bg-gradient-to-br from-purple-500 to-blue-500 text-white text-xl font-bold">
                      {profileData.name ? profileData.name[0] : "А"}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <h1 className="font-geist font-bold text-3xl text-white mb-1">
                      Добро пожаловать, {profileData.name || "Пользователь"}!
                    </h1>
                    <p className="text-white/70 font-manrope">Управляйте своим аккаунтом и настройками</p>
                  </div>
                </div>
              </div>

              {/* Stats Cards */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-white/70 text-sm font-manrope">Обработано файлов</p>
                        <p className="text-2xl font-bold font-geist">{stats.total_files}</p>
                      </div>
                      <div className="w-12 h-12 bg-purple-500/20 rounded-lg flex items-center justify-center">
                        <FileImage className="w-6 h-6 text-purple-400" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
                <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-white/70 text-sm font-manrope">Успешно обработано</p>
                        <p className="text-2xl font-bold font-geist">{stats.total_processed}</p>
                      </div>
                      <div className="w-12 h-12 bg-blue-500/20 rounded-lg flex items-center justify-center">
                        <Video className="w-6 h-6 text-blue-400" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
                <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-white/70 text-sm font-manrope">Неудачных обработок</p>
                        <p className="text-2xl font-bold font-geist">{stats.total_failed}</p>
                      </div>
                      <div className="w-12 h-12 bg-gradient-to-r from-orange-500 to-red-500 rounded-lg flex items-center justify-center mx-auto mb-4">
                        <CircleX className="w-6 h-6 text-white" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
                <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-white/70 text-sm font-manrope">Общий объем</p>
                        <p className="text-2xl font-bold font-geist">{formatTotalSizeMB(stats.total_size)} МБ</p>
                      </div>
                      <div className="w-12 h-12 bg-yellow-500/20 rounded-lg flex items-center justify-center">
                        <Download className="w-6 h-6 text-yellow-400" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Main Dashboard Grid */}
              <div className="grid lg:grid-cols-3 gap-8">
                {/* Left Column */}
                <div className="lg:col-span-2 space-y-6">
                  {/* Account Information */}
                  <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2 font-geist">
                        <User className="w-5 h-5" />
                        Информация об аккаунте
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      {editing ? (
                        <div className="grid md:grid-cols-2 gap-4">
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Имя</label>
                            <Input
                              value={profileData.name}
                              onChange={(e) => setProfileData({ ...profileData, name: e.target.value })}
                              className="bg-white/5 border-white/10 text-white placeholder:text-white/70"
                              placeholder="Введите имя"
                            />
                            {errors.name && <p className="text-red-400 text-sm">{errors.name}</p>}
                          </div>
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Email</label>
                            <Input
                              value={profileData.email}
                              onChange={(e) => setProfileData({ ...profileData, email: e.target.value })}
                              className="bg-white/5 border-white/10 text-white placeholder:text-white/70"
                              placeholder="Введите email"
                              type="email"
                            />
                            {errors.email && <p className="text-red-400 text-sm">{errors.email}</p>}
                          </div>
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Новый пароль (опционально)</label>
                            <Input
                              value={profileData.password}
                              onChange={(e) => setProfileData({ ...profileData, password: e.target.value })}
                              className="bg-white/5 border-white/10 text-white placeholder:text-white/70"
                              placeholder="Введите новый пароль"
                              type="password"
                            />
                            {errors.password && <p className="text-red-400 text-sm">{errors.password}</p>}
                          </div>
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Дата регистрации</label>
                            <div className="flex items-center gap-2 p-3 bg-white/5 rounded-lg border border-white/10">
                              <Calendar className="w-4 h-4 text-white/50" />
                              <span className="font-manrope">{formatDate(stats.created_at)}</span>
                            </div>
                          </div>
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Последнее обновление статистики</label>
                            <div className="flex items-center gap-2 p-3 bg-white/5 rounded-lg border border-white/10">
                              <Calendar className="w-4 h-4 text-white/50" />
                              <span className="font-manrope">{formatDate(stats.last_stats_update)}</span>
                            </div>
                          </div>
                          <div className="flex gap-3 col-span-2">
                            <Button
                              onClick={handleUpdateProfile}
                              className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                            >
                              Сохранить
                            </Button>
                            <Button
                              variant="outline"
                              className="bg-white/5 border-white/20 text-white hover:bg-white/10"
                              onClick={() => {
                                setEditing(false);
                                setErrors({ name: "", email: "", password: "" });
                                setProfileData(initialProfileData);
                              }}
                            >
                              Отмена
                            </Button>
                          </div>
                        </div>
                      ) : (
                        <div className="grid md:grid-cols-2 gap-4">
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Имя</label>
                            <div className="flex items-center justify-between p-3 bg-white/5 rounded-lg border border-white/10">
                              <span className="font-manrope">{profileData.name || "Не указано"}</span>
                              <Edit3
                                className="w-4 h-4 text-white/50 cursor-pointer hover:text-white"
                                onClick={() => setEditing(true)}
                              />
                            </div>
                          </div>
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Email</label>
                            <div className="flex items-center justify-between p-3 bg-white/5 rounded-lg border border-white/10">
                              <span className="font-manrope">{profileData.email || "Не указано"}</span>
                              <Edit3
                                className="w-4 h-4 text-white/50 cursor-pointer hover:text-white"
                                onClick={() => setEditing(true)}
                              />
                            </div>
                          </div>
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Дата регистрации</label>
                            <div className="flex items-center gap-2 p-3 bg-white/5 rounded-lg border border-white/10">
                              <Calendar className="w-4 h-4 text-white/50" />
                              <span className="font-manrope">{formatDate(stats.created_at)}</span>
                            </div>
                          </div>
                          <div className="space-y-2">
                            <label className="text-sm text-white/70 font-manrope">Последнее обновление статистики</label>
                            <div className="flex items-center gap-2 p-3 bg-white/5 rounded-lg border border-white/10">
                              <Calendar className="w-4 h-4 text-white/50" />
                              <span className="font-manrope">{formatDate(stats.last_stats_update)}</span>
                            </div>
                          </div>
                        </div>
                      )}
                    </CardContent>
                  </Card>

                  {/* Subscription Status (Placeholder) */}
                  <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2 font-geist">
                        <Crown className="w-5 h-5 text-yellow-400" />
                        Подписка Pro
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="font-manrope text-lg">Безлимитная обработка</p>
                        </div>
                        <Badge className="bg-green-500/20 text-green-400 border-green-400/30">Активна</Badge>
                      </div>
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-white/70 font-manrope">Использовано в этом месяце</span>
                          <span className="font-manrope">{stats.total_files} из ∞</span>
                        </div>
                        <Progress value={(stats.total_files / 300) * 100} className="h-2" />
                      </div>
                    </CardContent>
                  </Card>
                </div>

                {/* Right Column */}
                <div className="space-y-6">
                  {/* Quick Actions */}
                  <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                    <CardHeader>
                      <CardTitle className="font-geist">Быстрые действия</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      <Button
                        asChild
                        className="w-full justify-start bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                      >
                        <Link href="/process">
                          <Camera className="w-4 h-4 mr-2" />
                          Обработать файл
                        </Link>
                      </Button>
                      <Button
                        asChild
                        variant="outline"
                        className="w-full justify-start bg-white/5 border-white/20 text-white hover:bg-white/10"
                      >
                        <Link href="/history">
                          <Clock className="w-4 h-4 mr-2" />
                          Посмотреть историю
                        </Link>
                      </Button>
                      
                    </CardContent>
                  </Card>

                  {/* Recent Activity (Placeholder) ДОБАВИТЬ ПОСЛЕДНИЕ ТРИ ФАЙЛА, КТО ТО УМНЫЙ СДЕЛАЕТ */}
                  <Card className="bg-white/5 backdrop-blur-lg border-white/10 text-white">
                    <CardHeader>
                      <CardTitle className="font-geist">Последняя активность</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      <p className="text-white/70 text-sm">Нет недавней активности</p>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </>
          ) : (
            <div className="text-center py-16">
              <div className="w-24 h-24 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-6">
                <User className="w-12 h-12 text-white/40" />
              </div>
              <h3 className="font-geist font-semibold text-xl text-white mb-2">Войдите, чтобы увидеть личный кабинет</h3>
              <p className="text-white/70 mb-6">Авторизуйтесь, чтобы управлять вашим аккаунтом</p>
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
      </div>
    </div>
  );
}