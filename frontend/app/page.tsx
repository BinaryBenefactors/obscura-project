"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { RegistrationModal } from "@/components/registration-modal";
import { LoginModal } from "@/components/login-modal";
import { Upload, Settings, Download, Camera, Video, Users, Eye, Target, Cloud } from "lucide-react";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useAuth } from "@/components/AuthContext";

export default function HomePage() {
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [showRegistrationModal, setShowRegistrationModal] = useState(false);
  const { user, logout, isAuthenticated } = useAuth();

  const handleSwitchToLogin = () => {
    setShowRegistrationModal(false);
    setTimeout(() => setShowLoginModal(true), 100);
  };

  const handleSwitchToRegister = () => {
    setShowLoginModal(false);
    setTimeout(() => setShowRegistrationModal(true), 100);
  };

  useEffect(() => {
    const createParticles = () => {
      const particlesContainer = document.getElementById("particles");
      if (!particlesContainer) return;

      for (let i = 0; i < 50; i++) {
        const particle = document.createElement("div");
        particle.className = "particle";
        particle.style.left = Math.random() * 100 + "%";
        particle.style.animationDelay = Math.random() * 15 + "s";
        particle.style.animationDuration = Math.random() * 10 + 10 + "s";
        particlesContainer.appendChild(particle);
      }
    };

    createParticles();
  }, []);

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="fixed top-0 w-full z-50 backdrop-blur-lg bg-black/10 border-b border-white/10 transition-all duration-300">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-br from-white to-purple-400 rounded-lg flex items-center justify-center">
                <Camera className="w-5 h-5 text-black" />
              </div>
              <span className="font-geist font-bold text-xl bg-gradient-to-r from-white to-purple-400 bg-clip-text text-transparent">
                Obscura
              </span>
            </div>
            <nav className="hidden md:flex items-center gap-6">
              <a href="#how-it-works" className="text-white/80 hover:text-white transition-colors relative group">
                Как работает
                <span className="absolute bottom-[-5px] left-0 w-0 h-0.5 bg-gradient-to-r from-purple-400 to-blue-500 transition-all duration-300 group-hover:w-full"></span>
              </a>
              <a href="#benefits" className="text-white/80 hover:text-white transition-colors relative group">
                Преимущества
                <span className="absolute bottom-[-5px] left-0 w-0 h-0.5 bg-gradient-to-r from-purple-400 to-blue-500 transition-all duration-300 group-hover:w-full"></span>
              </a>
              <a href="#demo" className="text-white/80 hover:text-white transition-colors relative group">
                Демо
                <span className="absolute bottom-[-5px] left-0 w-0 h-0.5 bg-gradient-to-r from-purple-400 to-blue-500 transition-all duration-300 group-hover:w-full"></span>
              </a>
              <a href="#technology" className="text-white/80 hover:text-white transition-colors relative group">
                Технологии
                <span className="absolute bottom-[-5px] left-0 w-0 h-0.5 bg-gradient-to-r from-purple-400 to-blue-500 transition-all duration-300 group-hover:w-full"></span>
              </a>
              {isAuthenticated && (
                <Link href="/history" className="text-white/80 hover:text-white transition-colors relative group">
                  История
                  <span className="absolute bottom-[-5px] left-0 w-0 h-0.5 bg-gradient-to-r from-purple-400 to-blue-500 transition-all duration-300 group-hover:w-full"></span>
                </Link>
              )}
            </nav>
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
                    <Button className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 shadow-lg hover:shadow-purple-500/25 transition-all duration-300 transform hover:-translate-y-1 relative overflow-hidden">
                      Регистрация
                    </Button>
                  </RegistrationModal>
                </>
              )}
            </div>
          </div>
        </div>
      </header>

      <section className="hero-section relative bg-black flex items-center justify-center min-h-screen">
        <div className="bg-animation absolute top-0 left-0 w-full h-full z-0"></div>
        <div className="particles absolute top-0 left-0 w-full h-full z-10" id="particles"></div>
        <div className="hero-content text-center max-w-4xl mx-auto px-6 relative z-20">
          <h1 className="font-geist font-black text-6xl lg:text-8xl xl:text-9xl leading-none mb-8 bg-gradient-to-r from-white via-purple-400 to-blue-500 bg-clip-text text-transparent animate-fade-in-up">
            Obscura
          </h1>
          <p className="font-manrope text-xl lg:text-2xl text-white/70 mb-12 leading-relaxed font-light animate-fade-in-up-delay-200">
            Размывайте конфиденциальные объекты
            <br /> автоматически - без ручной работы
          </p>
          <div className="hero-actions flex flex-col sm:flex-row gap-6 justify-center animate-fade-in-up-delay-400">
            <Button
              size="lg"
              className="primary-btn font-manrope text-lg px-12 py-6 bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 shadow-xl hover:shadow-purple-500/30 transition-all duration-300 transform hover:-translate-y-1 relative overflow-hidden"
              asChild
            >
              <Link href="/process">Попробовать сейчас</Link>
            </Button>
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section id="how-it-works" className="py-20 bg-muted/30">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="font-geist font-bold text-3xl lg:text-4xl text-foreground mb-4">Как это работает</h2>
            <p className="font-manrope text-xl text-muted-foreground max-w-2xl mx-auto">
              Четкий алгоритм из 4 шагов с упором на поддержку профессиональных форматов
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            <Card className="text-center p-6">
              <CardContent className="pt-6">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Upload className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-lg mb-2">1. Загрузите фото или видео</h3>
                <p className="font-manrope text-muted-foreground text-sm">Доступны все форматы: JPEG, PNG, MP4, MOV</p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="pt-6">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Camera className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-lg mb-2">2. Выделите объекты для размытия</h3>
                <p className="font-manrope text-muted-foreground text-sm">
                  ИИ автоматически находит лица, автомобили, документы
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="pt-6">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Settings className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-lg mb-2">3. Настройте уровень защиты</h3>
                <p className="font-manrope text-muted-foreground text-sm">
                  Выберите интенсивность и тип эффекта
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="pt-6">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Download className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-lg mb-2">4. Скачайте результат</h3>
                <p className="font-manrope text-muted-foreground text-sm">
                  Получите обработанный файл без водяных знаков
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* Benefits Section */}
      <section id="benefits" className="py-20 bg-background">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="font-geist font-bold text-3xl lg:text-4xl text-foreground mb-4">
              Почему стоит выбрать Obscura
            </h2>
            <p className="font-manrope text-xl text-muted-foreground max-w-2xl mx-auto">
              Преимущества для фотографов, видеографов и всех, кто ценит конфиденциальность
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            <Card className="p-6">
              <CardContent className="pt-6">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Video className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-lg mb-2">Поддержка видео</h3>
                <p className="font-manrope text-muted-foreground text-sm">
                  Размывайте конфиденциальные объекты в видео с той же легкостью, что и на фото
                </p>
              </CardContent>
            </Card>

            <Card className="p-6">
              <CardContent className="pt-6">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Users className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-lg mb-2">Масштабируемость</h3>
                <p className="font-manrope text-muted-foreground text-sm">
                  Обрабатывайте сотни файлов одновременно с высокой скоростью
                </p>
              </CardContent>
            </Card>

            <Card className="p-6">
              <CardContent className="pt-6">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Settings className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-lg mb-2">Гибкие настройки</h3>
                <p className="font-manrope text-muted-foreground text-sm">
                  Выбирайте, что и как размывать: лица, номера, документы
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* Interactive Demo Zone */}
      <section id="demo" className="py-20 bg-muted/30">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="font-geist font-bold text-3xl lg:text-4xl text-foreground mb-4">
              Убедитесь в качестве сами
            </h2>
            <p className="font-manrope text-xl text-muted-foreground max-w-3xl mx-auto">
              Посмотрите, как AI точно находит и скрывает конфиденциальные данные. Просто наведите курсор на объекты на
              изображении ниже.
            </p>
          </div>

          <div className="max-w-4xl mx-auto">
            <div className="relative bg-card rounded-xl p-8 mb-8">
              <div className="relative mx-auto max-w-2xl">
                <img
                  src="/street-scene-with-people-and-cars-for-ai-detection.png"
                  alt="Демо изображение для интерактивной демонстрации"
                  className="w-full rounded-lg"
                />

                {/* Interactive overlay areas */}
                <div className="absolute inset-0 rounded-lg">
                  {/* Face detection area */}
                  <div
                    className="absolute top-[20%] left-[30%] w-16 h-20 border-2 border-transparent hover:border-green-400 hover:bg-green-400/20 rounded cursor-pointer transition-all duration-200 group"
                    title="Лицо"
                  >
                    <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-green-600 text-white px-2 py-1 rounded text-xs opacity-0 group-hover:opacity-100 transition-opacity">
                      Лицо
                    </div>
                  </div>

                  {/* License plate area */}
                  <div
                    className="absolute bottom-[25%] right-[20%] w-20 h-8 border-2 border-transparent hover:border-green-400 hover:bg-green-400/20 rounded cursor-pointer transition-all duration-200 group"
                    title="Номерной знак"
                  >
                    <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-green-600 text-white px-2 py-1 rounded text-xs opacity-0 group-hover:opacity-100 transition-opacity">
                      Номерной знак
                    </div>
                  </div>

                  {/* Another face area */}
                  <div
                    className="absolute top-[35%] right-[35%] w-14 h-18 border-2 border-transparent hover:border-green-400 hover:bg-green-400/20 rounded cursor-pointer transition-all duration-200 group"
                    title="Лицо"
                  >
                    <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-green-600 text-white px-2 py-1 rounded text-xs opacity-0 group-hover:opacity-100 transition-opacity">
                      Лицо
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Demo Controls */}
            <Card className="p-6">
              <div className="space-y-6">
                <div>
                  <label className="font-manrope font-medium text-sm mb-3 block">Интенсивность эффекта</label>
                  <div className="relative">
                    <input
                      type="range"
                      min="1"
                      max="10"
                      defaultValue="5"
                      className="w-full h-2 rounded-lg slider-with-fill"
                      style={{
                        background: "linear-gradient(to right, #2563eb 0%, #2563eb 45%, #d1d5db 0%, #d1d5db 100%)",
                      }}
                      onChange={(e) => {
                        const value = ((e.target.value - e.target.min) / (e.target.max - e.target.min)) * 100
                        e.target.style.background = `linear-gradient(to right, #2563eb 0%, #2563eb ${value}%, #d1d5db ${value}%, #d1d5db 100%)`
                      }}
                      onInput={(e) => {
                        const value = ((e.target.value - e.target.min) / (e.target.max - e.target.min)) * 100
                        e.target.style.background = `linear-gradient(to right, #2563eb 0%, #2563eb ${value}%, #d1d5db ${value}%, #d1d5db 100%)`
                      }}
                      onLoad={(e) => {
                        const value = ((e.target.value - e.target.min) / (e.target.max - e.target.min)) * 100
                        e.target.style.background = `linear-gradient(to right, #2563eb 0%, #2563eb ${value}%, #d1d5db ${value}%, #d1d5db 100%)`
                      }}
                    />
                  </div>
                  <div className="flex justify-between text-xs text-muted-foreground mt-1">
                    <span>Слабое</span>
                    <span>Сильное</span>
                  </div>
                </div>

                <div>
                  <label className="font-manrope font-medium text-sm mb-3 block">Тип эффекта</label>
                  <div className="flex gap-3">
                    <Button variant="outline" size="sm" className="font-manrope bg-transparent">
                      Размытие
                    </Button>
                    <Button variant="outline" size="sm" className="font-manrope bg-transparent">
                      Пикселизация
                    </Button>
                    <Button variant="outline" size="sm" className="font-manrope bg-transparent">
                      Маска
                    </Button>
                  </div>
                </div>

                <div className="pt-4 border-t">
                  <Button size="lg" className="w-full font-manrope" asChild>
                    <Link href="/process">Так обработать мое фото</Link>
                  </Button>
                </div>
              </div>
            </Card>
          </div>
        </div>
      </section>

      {/* Technologies Under the Hood */}
      <section id="technology" className="py-20">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="font-geist font-bold text-3xl lg:text-4xl text-foreground mb-4">
              Машинное обучение на страже вашей конфиденциальности
            </h2>
            <p className="font-manrope text-xl text-muted-foreground max-w-3xl mx-auto">
              Obscura — это не просто фильтр, а сложная система компьютерного зрения.
            </p>
          </div>

          <div className="grid lg:grid-cols-3 gap-8">
            <Card className="p-8 text-center">
              <CardContent className="pt-0">
                <div className="w-16 h-16 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-6">
                  <Eye className="w-8 h-8 text-primary" />
                </div>
                <h3 className="font-geist font-semibold text-xl mb-4">Детекция объектов</h3>
                <p className="font-manrope text-muted-foreground leading-relaxed">
                  Используем дообученную модель YOLO для молниеносного и точного распознавания лиц, номерных знаков и
                  других объектов в кадре.
                </p>
              </CardContent>
            </Card>

            <Card className="p-8 text-center">
              <CardContent className="pt-0">
                <div className="w-16 h-16 bg-accent/10 rounded-lg flex items-center justify-center mx-auto mb-6">
                  <Target className="w-8 h-8 text-accent" />
                </div>
                <h3 className="font-geist font-semibold text-xl mb-4">Высокая точность</h3>
                <p className="font-manrope text-muted-foreground leading-relaxed">
                  Наши алгоритмы обучены на разнообразных данных, чтобы работать одинаково хорошо при любом освещении и
                  ракурсе.
                </p>
              </CardContent>
            </Card>

            <Card className="p-8 text-center">
              <CardContent className="pt-0">
                <div className="w-16 h-16 bg-secondary/10 rounded-lg flex items-center justify-center mx-auto mb-6">
                  <Cloud className="w-8 h-8 text-secondary-foreground" />
                </div>
                <h3 className="font-geist font-semibold text-xl mb-4">Мощь облаков</h3>
                <p className="font-manrope text-muted-foreground leading-relaxed">
                  Вся обработка происходит на наших серверах. Вам не нужно мощное железо для работы с 4K и 8K видео.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-primary/5">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="font-geist font-bold text-3xl lg:text-4xl text-foreground mb-6">Готовы начать?</h2>
          <p className="font-manrope text-xl text-muted-foreground mb-8">
            Попробуйте AI Blur бесплатно и убедитесь в качестве автоматического размытия
          </p>
          <Button size="lg" className="font-manrope text-lg px-8 py-6" asChild>
            <Link href="/process">Попробовать бесплатно</Link>
          </Button>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border bg-card/50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="flex items-center gap-2 mb-4 md:mb-0">
              <div className="w-8 h-8 bg-gradient-to-br from-white to-purple-400 rounded-lg flex items-center justify-center">
                <Camera className="w-5 h-5 text-black" />
              </div>
              <span className="font-geist font-semibold text-xl bg-gradient-to-r from-white to-purple-400 bg-clip-text text-transparent">
                Obscura
              </span>
            </div>
            <div className="flex gap-6 text-sm text-muted-foreground">
              <a href="https://github.com/ai-blur/obscura" className="hover:text-foreground transition-colors">
                GitHub Repository
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}