"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card";
import { RegistrationModal } from "@/components/registration-modal"
import { LoginModal } from "@/components/login-modal"
import Link from "next/link"
import { useEffect, useState } from "react"
import { User, Settings, LogOut, ChevronDown, ChevronUp } from "lucide-react"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { useAuth } from "@/components/AuthContext";
import CameraIcon from "@/components/ui/camera-icon";
import { DemoCanvas } from "@/components/demo-canvas"

export default function HomePage() {
  const [showLoginModal, setShowLoginModal] = useState(false)
  const [showRegistrationModal, setShowRegistrationModal] = useState(false)
  const { user, logout, isAuthenticated } = useAuth();
  const [effectIntensity, setEffectIntensity] = useState(5);
  const [open, setOpen] = useState(false)

  const handleSwitchToLogin = () => {
    setShowRegistrationModal(false)
    setTimeout(() => setShowLoginModal(true), 100)
  }

  const handleSwitchToRegister = () => {
    setShowLoginModal(false)
    setTimeout(() => setShowRegistrationModal(true), 100)
  }

  useEffect(() => {
    // Логика курсора
    const cursor = document.querySelector('.cursor');
    const cursorFollower = document.querySelector('.cursor-follower');
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
      cursor?.classList.add('active');
    };

    const handleMouseUp = () => {
      cursor?.classList.remove('active');
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mousedown', handleMouseDown);
    document.addEventListener('mouseup', handleMouseUp);

    // Создание частиц
    const createParticles = () => {
      const particles = document.getElementById('particles');
      if (!particles) return;

      const particleCount = 30;

      for (let i = 0; i < particleCount; i++) {
        const particle = document.createElement('div');
        particle.className = 'particle';
        particle.style.width = Math.random() * 4 + 1 + 'px';
        particle.style.height = particle.style.width;
        particle.style.left = Math.random() * 100 + '%';
        particle.style.animationDuration = Math.random() * 20 + 10 + 's';
        particle.style.animationDelay = Math.random() * 20 + 's';
        particle.style.animation = `particle-up ${particle.style.animationDuration} linear infinite`;
        particles.appendChild(particle);
      }
    };

    createParticles();

    // Логика прокрутки заголовка
    const handleScroll = () => {
      const header = document.getElementById('header');
      if (header) {
        if (window.scrollY > 50) {
          header.classList.add('scrolled');
        } else {
          header.classList.remove('scrolled');
        }
      }
    };

    window.addEventListener('scroll', handleScroll);

    // Логика переключателей
    document.querySelectorAll('.toggle-group').forEach((group) => {
      const buttons = group.querySelectorAll('.toggle-btn');
      buttons.forEach((btn) => {
        btn.addEventListener('click', () => {
          buttons.forEach((b) => b.classList.remove('active'));
          btn.classList.add('active');
        });
      });
    });

    // Инициализация видимости .feature-card
    const cards = document.querySelectorAll('.feature-card');

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add('visible');
          }
        });
      },
      {
        threshold: 0.1, // Анимация начинается, когда 10% элемента видно
      }
    );

    cards.forEach((card) => {
      observer.observe(card);
    });

    // Показ курсора
    if (cursor) cursor.style.opacity = '1';
    if (cursorFollower) cursorFollower.style.opacity = '1';

    // Очистка всех слушателей
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mousedown', handleMouseDown);
      document.removeEventListener('mouseup', handleMouseUp);
      window.removeEventListener('scroll', handleScroll);
      observer.disconnect(); // Очистка Intersection Observer
    };
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <div className="noise"></div>
      <div className="cursor"></div>
      <div className="cursor-follower"></div>

      <header id="header" className="header">
        <div className="header-container">
          <Link href="/" className="logo">
            <div className="logo-icon">
              <CameraIcon />
            </div>
            <span className="logo-text">Obscura</span>
          </Link>

          <nav className="nav">
            <a href="#features" className="nav-link">
              Возможности
            </a>
            <a href="#demo" className="nav-link">
              Демо
            </a>
            <Link href="/history" className="nav-link">
              История
            </Link>
          </nav>

          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <DropdownMenu open={open} onOpenChange={setOpen}>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="outline"
                    className="flex items-center gap-2 font-manrope text-white bg-white/10 hover:bg-white/20 border-0"
                  >
                    <User className="h-4 w-4" />
                    {user?.name || user?.email || "Пользователь"}
                    {open ? (
                      <ChevronUp className="h-4 w-4 ml-1" />
                    ) : (
                      <ChevronDown className="h-4 w-4 ml-1" />
                    )}
                  </Button>
                </DropdownMenuTrigger>

                <DropdownMenuContent align="end" className="w-56 h-48 z-2000">
                  <div>
                    <DropdownMenuLabel>
                      <div className="flex flex-col">
                        <span className="font-semibold">{user?.name || "Пользователь"}</span>
                        <span className="text-sm text-[#8c939f]">{user?.email}</span>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                  </div>

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
                  <Button className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 shadow-lg hover:shadow-purple-500/25 transition-all duration-300 transform hover:-translate-y-1">
                    Регистрация
                  </Button>
                </RegistrationModal>
              </>
            )}
          </div>
        </div>
      </header>

      <section className="hero">
        <div className="hero-bg">
          <div className="floating-shapes">
            <div className="shape shape-1"></div>
            <div className="shape shape-2"></div>
            <div className="shape shape-3"></div>
          </div>
          <div className="grid-overlay"></div>
          <div className="particles" id="particles"></div>
        </div>

        <div className="hero-content">
          <div className="hero-badge">
            <div className="hero-badge-dot"></div>
            <span>Powered by Advanced AI</span>
          </div>

          <h1 className="hero-title">
            <span className="title-main">Obscura</span>
          </h1>

          <p className="hero-subtitle">
            Интеллектуальная система защиты конфиденциальности.
            <br />
            <span>Автоматическое обнаружение</span> и <span>мгновенное скрытие</span>
            личных данных в медиаконтенте
          </p>

          <div className="hero-actions">
            <Link href="/process" className="btn-primary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="white">
                <polygon points="5,3 19,12 5,21" />
              </svg>
              Начать сейчас
            </Link>
            <a href="#demo" className="btn-secondary">
              <span>Посмотреть демо</span>
            </a>
          </div>
        </div>

        <div className="scroll-indicator">
          <svg viewBox="0 0 24 24">
            <path d="m7 10 5 5 5-5" />
            <path d="m7 4 5 5 5-5" />
          </svg>
        </div>
      </section>

      <section className="features" id="features">
        <div className="container">
          <div className="section-header">
            <span className="section-tag">Возможности</span>
            <h2 className="section-title">Технологии будущего</h2>
            <p className="section-subtitle">Передовые алгоритмы машинного обучения для максимальной защиты</p>
          </div>

          <div className="features-grid">
            <div className="feature-card">
              <div className="feature-icon">
                <svg viewBox="0 0 24 24">
                  <path d="M9.5,3A6.5,6.5 0 0,1 16,9.5C16,11.11 15.41,12.59 14.44,13.73L14.71,14H15.5L20.5,19L19,20.5L14,15.5V14.71L13.73,14.44C12.59,15.41 11.11,16 9.5,16A6.5,6.5 0 0,1 3,9.5A6.5,6.5 0 0,1 9.5,3M9.5,5C7,5 5,7 5,9.5C5,12 7,14 9.5,14C12,14 14,12 14,9.5C14,7 12,5 9.5,5Z" />
                </svg>
              </div>
              <h3 className="feature-title">Умное распознавание</h3>
              <p className="feature-description">
                AI автоматически находит лица, номерные знаки, документы и другие конфиденциальные объекты с точностью
                99.9%
              </p>
            </div>

            <div className="feature-card">
              <div className="feature-icon">
                <svg viewBox="0 0 24 24">
                  <path d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2M12,4A8,8 0 0,1 20,12C20,14.4 19,16.5 17.3,18C15.9,16.7 14,16 12,16C10,16 8.2,16.7 6.7,18C5,16.5 4,14.4 4,12A8,8 0 0,1 12,4M14,5.89C13.62,5.9 13.26,6.15 13.1,6.54L11.81,9.77L11.71,10C11,10.13 10.41,10.6 10.14,11.26C9.73,12.29 10.23,13.45 11.26,13.86C12.29,14.27 13.45,13.77 13.86,12.74C14.12,12.08 14,11.32 13.57,10.76L13.67,10.5L14.96,7.29L14.97,7.26C15.17,6.75 14.92,6.17 14.41,5.96C14.28,5.91 14.15,5.89 14,5.89M10,6A1,1 0 0,0 9,7A1,1 0 0,0 10,8A1,1 0 0,0 11,7A1,1 0 0,0 10,6M7,9A1,1 0 0,0 6,10A1,1 0 0,0 7,11A1,1 0 0,0 8,10A1,1 0 0,0 7,9Z" />
                </svg>
              </div>
              <h3 className="feature-title">Обработка в реальном времени</h3>
              <p className="feature-description">
                Мгновенная обработка видео 4K со скоростью 60 FPS без потери качества и задержек
              </p>
            </div>

            <div className="feature-card">
              <div className="feature-icon">
                <svg viewBox="0 0 24 24">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
              </div>
              <h3 className="feature-title">Полная конфиденциальность</h3>
              <p className="feature-description">
                Вся обработка происходит локально. Ваши данные никогда не покидают ваше устройство
              </p>
            </div>

            <div className="feature-card">
              <div className="feature-icon">
                <svg viewBox="0 0 24 24">
                  <path d="M12,18.17L8.83,15L7.42,16.41L12,21L16.59,16.41L15.17,15M12,5.83L15.17,9L16.58,7.59L12,3L7.41,7.59L8.83,9L12,5.83Z" />
                </svg>
              </div>
              <h3 className="feature-title">Гибкие настройки</h3>
              <p className="feature-description">
                Выбирайте что скрывать: лица, номера, документы. Настраивайте тип и интенсивность эффектов
              </p>
            </div>

            <div className="feature-card">
              <div className="feature-icon">
                <svg viewBox="0 0 24 24">
                  <path d="M3,13H5V11H3V13M3,17H5V15H3V17M3,9H5V7H3V9M7,13H21V11H7V13M7,17H21V15H7V17M7,7V9H21V7H7Z" />
                </svg>
              </div>
              <h3 className="feature-title">Пакетная обработка</h3>
              <p className="feature-description">
                Обрабатывайте тысячи файлов одновременно с автоматической оптимизацией ресурсов
              </p>
            </div>

            <div className="feature-card">
              <div className="feature-icon">
                <svg viewBox="0 0 24 24">
                  <path d="M12,2A3,3 0 0,1 15,5V11A3,3 0 0,1 12,14A3,3 0 0,1 9,11V5A3,3 0 0,1 12,2M19,11C19,14.53 16.39,17.44 13,17.93V21H11V17.93C7.61,17.44 5,14.53 5,11H7A5,5 0 0,0 12,16A5,5 0 0,0 17,11H19Z" />
                </svg>
              </div>
              <h3 className="feature-title">API для разработчиков</h3>
              <p className="feature-description">
                Простая интеграция в любые проекты через REST API с подробной документацией
              </p>
            </div>
          </div>
        </div>
      </section>

      <DemoCanvas />

      <section className="cta">
        <div className="container">
          <div className="cta-content">
            <h2 className="cta-title">Готовы защитить свою конфиденциальность?</h2>
            <p className="cta-text">Присоединяйтесь к тысячам пользователей, которые уже доверяют Obscura</p>
            <Link href="/process" className="btn-primary">
              Попробовать бесплатно
            </Link>
          </div>
        </div>
      </section>

      <footer className="footer">
        <div className="container">
          <div className="footer-content">
            <div className="footer-brand">
              <div className="footer-logo">
                <div className="logo-icon">
                  <CameraIcon />
                </div>
                <span className="logo-text">Obscura</span>
              </div>
              <p className="footer-description">
                Передовая система автоматической защиты конфиденциальности в медиаконтенте на основе искусственного
                интеллекта
              </p>
            </div>
            <div className="footer-section">
              <h4>Ссылки</h4>
              <ul className="footer-links">
                <li>
                  <a href="#features">Возможности</a>
                </li>
                <li>
                  <a href="#demo">Демонстрация</a>
                </li>
                <li>
                  <a href="https://github.com/BinaryBenefactors/obscura-project">GitHub</a>
                </li>
              </ul>
            </div>
          </div>

          <div className="footer-bottom">
            <div className="footer-copyright">© 2025 Obscura. Все права защищены.</div>
          </div>
        </div>
      </footer>
    </div>
  )
}
