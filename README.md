# Проект «Obscura»
---
## Введение

Проект **«Obscura»** был создан с целью упрощения процесса заблюривания объектов на изображениях и видео, что позволяет пользователям легко скрывать конфиденциальную информацию или личные данные, а также маскировать нежелательные объекты в реальном времени. Это приложение ориентировано на пользователей, которым важно защищать свою конфиденциальность и безопасности данных.

## Команда

### Наставник:
- Семенов Владимир Евгеньевич

### Члены команды:
**Капитан**: Ягольник Даниил Сергеевич

- **Батычков Вячеслав Геннадьевич**, КТбо1-7, Backend developer, QA engineer.
- **Долбин Матвей Сергеевич**,       КТ6о1-7, Frontend developer, UX/UI designer.
- **Лавров Даниил Эдуардович**,      КТбо1-7, Frontend developer, ML-engineer.
- **Одинцов Дмитрий Максимович**,    КТбо1-7, ML-engineer.
- **Скубриев Роман Владимирович**,   КТсo1-4, DevOps-инженер, Project Manager.
- **Ягольник Даниил Сергеевич**,     КТбо1-7, Fullstack developer, Content Manager.

## О проекте

Проект **«Obscura»** — это веб-приложение, предназначенное для заблюривания объектов на фото/видео. Приложение позволяет пользователям выбирать объекты для блюра, а также настраивать степень и тип блюра и многое другое. Пользователи могут отправлять свои фотографии или видео для обработки и просматривать историю отправок.

### Основные функции:
- **Заблюривание объектов**: Пользователи могут заблюрить различные объекты на изображениях и видео.
- **Типы объектов**: Возможность выбора типа объекта для блюра.
- **Настройка блюра**: Выбор степени и типа блюра (размытие, пикселизация и т.д.).
- **Регистрация пользователей**: Пользователи могут зарегистрироваться для просмотра истории отправок и дополнительных возможностей.

### Коммерциализация:
Для использования приложения предусмотрена система токенов:
- Незарегистрированные пользователи имеют ограничения по выбору типа объектов, количеству обработок и качеству.
- Для зарегистрированных пользователей планируется расширение функционала (в том числе возможность дополнительных настроек обработки).

## Стек технологий

### Frontend:
- **React** — для создания пользовательского интерфейса.
- **TypeScript** — основной язык программирования.
- **HTML, CSS** — для оформления и разметки страниц.

### Backend:
- **Golang** — для серверной логики, обработки запросов и взаимодействия с базой данных.
- **PostgreSQL** — база данных для хранения информации о пользователях, объектах и истории отправок.

### Искусственный интеллект:
- **Python** — основной язык для разработки ML-моделей.
- **YOLO (You Only Look Once)** — модель для детекции объектов на изображениях и видео.

## Разработка

Проект состоит из нескольких компонентов, взаимодействующих друг с другом:

1. **Frontend** отвечает за отображение интерфейса и взаимодействие с пользователем.
2. **Backend** обрабатывает запросы, работает с базой данных и управляет логикой приложения.
3. **Искусственный интеллект** используется для детекции объектов на изображениях и видео с целью их последующего блюра.
4. **DevOps** отвечает за развертывание и поддержание приложения в рабочем состоянии.

## Структуры проекта

```
├── .git/ 🚫 (auto-hidden)
├── backend/
│   ├── backend-app/
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── docs/
│   │   │   ├── docs.go
│   │   │   ├── swagger.json
│   │   │   └── swagger.yaml
│   │   ├── internal/
│   │   │   ├── config.go
│   │   │   ├── database.go
│   │   │   ├── file_cleaner.go
│   │   │   ├── models.go
│   │   │   ├── rate_limiter.go
│   │   │   ├── server.go
│   │   │   └── validator.go
│   │   ├── pkg/
│   │   │   └── logger/
│   │   │       └── logger.go
│   │   ├── uploads/
│   │   ├── .env 🚫 (auto-hidden)
│   │   ├── .env.example
│   │   ├── .gitignore
│   │   ├── Dockerfile
│   │   ├── README.md
│   │   ├── app.log 🚫 (auto-hidden)
│   │   ├── go.mod
│   │   └── go.sum
│   └── ml/
│       ├── .env/ 🚫 (auto-hidden)
│       ├── app/
│       │   ├── __pycache__/ 🚫 (auto-hidden)
│       │   ├── ml/
│       │   │   ├── tools/
│       │   │   │   ├── __pycache__/ 🚫 (auto-hidden)
│       │   │   │   ├── model.py
│       │   │   │   ├── object_detector.py
│       │   │   │   └── write_box.py
│       │   │   └── ml_executor.py
│       │   ├── routers/
│       │   │   ├── __pycache__/ 🚫 (auto-hidden)
│       │   │   └── video.py
│       │   ├── schemas/
│       │   │   ├── __pycache__/ 🚫 (auto-hidden)
│       │   │   └── uploadfile.py
│       │   ├── tools/
│       │   │   ├── __pycache__/ 🚫 (auto-hidden)
│       │   │   └── generate_name_file.py
│       │   ├── main.py
│       │   └── requirements.txt
│       ├── media/
│       │   └── diplom.jpg
│       ├── .gitignore
│       ├── Dockerfile
│       ├── LICENSE
│       ├── api_test.py
│       └── load_model.py
├── docs/
│   ├── CONTRIBUTING.md
│   ├── deployment-strategy.md
│   └── github-secrets.md
├── frontend/
│   ├── app/
│   │   ├── dashboard/
│   │   │   └── page.tsx
│   │   ├── history/
│   │   │   └── page.tsx
│   │   ├── process/
│   │   │   └── page.tsx
│   │   ├── globals.css
│   │   ├── layout.tsx
│   │   └── page.tsx
│   ├── components/
│   │   ├── ui/
│   │   │   ├── accordion.tsx
│   │   │   ├── alert-dialog.tsx
│   │   │   ├── alert.tsx
│   │   │   ├── aspect-ratio.tsx
│   │   │   ├── avatar.tsx
│   │   │   ├── badge.tsx
│   │   │   ├── breadcrumb.tsx
│   │   │   ├── button.tsx
│   │   │   ├── calendar.tsx
│   │   │   ├── camera-icon.tsx
│   │   │   ├── card.tsx
│   │   │   ├── carousel.tsx
│   │   │   ├── chart.tsx
│   │   │   ├── checkbox.tsx
│   │   │   ├── collapsible.tsx
│   │   │   ├── command.tsx
│   │   │   ├── context-menu.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── drawer.tsx
│   │   │   ├── dropdown-menu.tsx
│   │   │   ├── form.tsx
│   │   │   ├── hover-card.tsx
│   │   │   ├── input-otp.tsx
│   │   │   ├── input.tsx
│   │   │   ├── label.tsx
│   │   │   ├── menubar.tsx
│   │   │   ├── navigation-menu.tsx
│   │   │   ├── pagination.tsx
│   │   │   ├── popover.tsx
│   │   │   ├── progress.tsx
│   │   │   ├── radio-group.tsx
│   │   │   ├── resizable.tsx
│   │   │   ├── scroll-area.tsx
│   │   │   ├── select.tsx
│   │   │   ├── separator.tsx
│   │   │   ├── sheet.tsx
│   │   │   ├── sidebar.tsx
│   │   │   ├── skeleton.tsx
│   │   │   ├── slider.tsx
│   │   │   ├── sonner.tsx
│   │   │   ├── switch.tsx
│   │   │   ├── table.tsx
│   │   │   ├── tabs.tsx
│   │   │   ├── textarea.tsx
│   │   │   ├── toast.tsx
│   │   │   ├── toaster.tsx
│   │   │   ├── toggle-group.tsx
│   │   │   ├── toggle.tsx
│   │   │   ├── tooltip.tsx
│   │   │   ├── use-mobile.tsx
│   │   │   └── use-toast.ts
│   │   ├── AuthContext.tsx
│   │   ├── login-modal.tsx
│   │   ├── registration-modal.tsx
│   │   └── theme-provider.tsx
│   ├── hooks/
│   │   ├── use-mobile.ts
│   │   └── use-toast.ts
│   ├── lib/
│   │   └── utils.ts
│   ├── public/
│   │   ├── placeholder-logo.png
│   │   ├── placeholder-logo.svg
│   │   ├── placeholder-user.jpg
│   │   ├── placeholder.jpg
│   │   ├── placeholder.svg
│   │   └── street-scene-with-people-and-cars-for-ai-detection.png
│   ├── styles/
│   │   └── globals.css
│   ├── .env 🚫 (auto-hidden)
│   ├── .env.example
│   ├── .gitignore
│   ├── .gitkeep
│   ├── components.json
│   ├── next.config.mjs
│   ├── package-lock.json
│   ├── package.json
│   ├── pnpm-lock.yaml
│   ├── postcss.config.mjs
│   └── tsconfig.json
├── nginx/
│   └── nginx.conf
├── uploads/
├── Dockerfile.frontend
├── Dockerfile.nginx
├── LICENSE
├── README.md
└── docker-compose.yml
```

## Планируемые улучшения

- Добавление новых типов блюра и масок.
- Оптимизация AI-алгоритмов для повышения точности детекции.
- Улучшение пользовательского опыта через адаптивный интерфейс и дополнительные настройки.
- Разработка и внедрение платных подписок для пользователей с расширенным функционалом.

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и развертывания. Подробную информацию о CI/CD пайплайне можно найти в [документации по CI/CD](.github/workflows/ci-cd.yml).

## Развёртывание

Подробная информация о развертывании проекта доступна в [документации по развертыванию](docs/deployment-strategy.md).

## Контрибьюции

Мы приветствуем контрибьюции! Для того чтобы предложить улучшения или исправления:

1. Форкните репозиторий.
2. Создайте новую ветку для вашего улучшения.
3. Напишите тесты и убедитесь, что они проходят.
4. Сделайте pull request.

Пожалуйста, следуйте [стилю кодирования](docs/CONTRIBUTING.md) и убедитесь, что ваш код соответствует стандартам проекта.

<div align="center">
    <h2> Список контрибьюторов </h2>
</div>
<div align="center">
    <a href="https://github.com/binarybenefactors/obscura-project/graphs/contributors">
      <img src="https://contrib.rocks/image?repo=binarybenefactors/obscura-project" />
    </a>
</div>

## Лицензия

Этот проект лицензирован под лицензией MIT. См. файл [LICENSE](LICENSE) для подробностей.
