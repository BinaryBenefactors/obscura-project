"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card";
import { RegistrationModal } from "@/components/registration-modal"
import { LoginModal } from "@/components/login-modal"
import Link from "next/link"
import { useEffect, useState } from "react"
import { useAuth } from "@/components/AuthContext";

export default function HomePage() {
  const [showLoginModal, setShowLoginModal] = useState(false)
  const [showRegistrationModal, setShowRegistrationModal] = useState(false)
  const { user, logout, isAuthenticated } = useAuth();
  const [effectIntensity, setEffectIntensity] = useState(5);

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
          <Link href="#" className="logo">
            <div className="logo-icon">
              <svg viewBox="0 0 24 24">
                <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z" />
                <circle cx="12" cy="13" r="4" />
              </svg>
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

      <section className="demo" id="demo">
        <div className="container">
          <div className="section-header">
            <span className="section-tag">Демонстрация</span>
            <h2 className="section-title">Увидьте магию в действии</h2>
            <p className="section-subtitle">Настройте параметры и посмотрите, как работает Obscura</p>
          </div>

          <div className="demo-content">
            <div className="demo-visual">
              <div className="demo-image-container">
                <img
                  src="https://images.unsplash.com/photo-1573164713988-8665fc963095?w=800&h=500&fit=crop"
                  alt="Demo"
                  className="demo-image"
                />
                <div className="demo-overlay"></div>
              </div>
            </div>

            <div className="demo-controls">
              <div className="control-group">
                <label className="control-label">Интенсивность размытия</label>
                <input type="range" min="0" max="100" defaultValue="50" className="slider" />
              </div>

              <div className="control-group">
                <label className="control-label">Тип эффекта</label>
                <div className="toggle-group">
                  <button className="toggle-btn active">Blur</button>
                  <button className="toggle-btn">Pixelate</button>
                  <button className="toggle-btn">Blackout</button>
                </div>
              </div>

              <div className="control-group">
                <label className="control-label">Объекты для скрытия</label>
                <div className="toggle-group">
                  <button className="toggle-btn active">Лица</button>
                  <button className="toggle-btn">Номера</button>
                  <button className="toggle-btn">Документы</button>
                </div>
              </div>

              <Link
                href="/process"
                className="btn-primary"
                style={{
                  width: "100%",
                  marginTop: "2rem",
                  display: "inline-flex",
                  alignItems: "center",
                  justifyContent: "center",
                  gap: "0.75rem",
                }}
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="white">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                </svg>
                Применить эффекты
              </Link>
            </div>
          </div>
        </div>
      </section>

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
                  <svg viewBox="0 0 24 24">
                    <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z" />
                    <circle cx="12" cy="13" r="4" />
                  </svg>
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
            <div className="footer-legal">
              <a href="/privacy">Конфиденциальность</a>
              <a href="/terms">Условия использования</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
