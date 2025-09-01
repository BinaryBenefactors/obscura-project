"use client";

import type React from "react";
import { useState, useCallback, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Slider } from "@/components/ui/slider";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { RegistrationModal } from "@/components/registration-modal";
import { LoginModal } from "@/components/login-modal";
import { Camera, Upload, Settings, Download, ArrowLeft, Search, Check, Trash } from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/components/AuthContext";

const API_BASE = "http://localhost:8080";

const RUS_TO_ENG_MAPPING: { [key: string]: string } = {
  "лицо": "face",
  "человек": "person",
  "велосипед": "bicycle",
  "автомобиль": "car",
  "мотоцикл": "motorcycle",
  "самолет": "airplane",
  "автобус": "bus",
  "поезд": "train",
  "грузовик": "truck",
  "лодка": "boat",
  "светофор": "traffic light",
  "пожарный гидрант": "fire hydrant",
  "знак стоп": "stop sign",
  "парковочный счетчик": "parking meter",
  "скамейка": "bench",
  "птица": "bird",
  "кот": "cat",
  "собака": "dog",
  "лошадь": "horse",
  "овца": "sheep",
  "корова": "cow",
  "слон": "elephant",
  "медведь": "bear",
  "зебра": "zebra",
  "жираф": "giraffe",
  "рюкзак": "backpack",
  "зонт": "umbrella",
  "сумка": "handbag",
  "галстук": "tie",
  "чемодан": "suitcase",
  "фрисби": "frisbee",
  "лыжи": "skis",
  "сноуборд": "snowboard",
  "спортивный мяч": "sports ball",
  "воздушный змей": "kite",
  "бейсбольная бита": "baseball bat",
  "бейсбольная перчатка": "baseball glove",
  "скейтборд": "skateboard",
  "доска для серфинга": "surfboard",
  "теннисная ракетка": "tennis racket",
  "бутылка": "bottle",
  "бокал для вина": "wine glass",
  "чашка": "cup",
  "вилка": "fork",
  "нож": "knife",
  "ложка": "spoon",
  "миска": "bowl",
  "банан": "banana",
  "яблоко": "apple",
  "бутерброд": "sandwich",
  "апельсин": "orange",
  "брокколи": "broccoli",
  "морковь": "carrot",
  "хот-дог": "hot dog",
  "пицца": "pizza",
  "пончик": "donut",
  "торт": "cake",
  "стул": "chair",
  "диван": "couch",
  "горшечное растение": "potted plant",
  "кровать": "bed",
  "обеденный стол": "dining table",
  "туалет": "toilet",
  "телевизор": "tv",
  "ноутбук": "laptop",
  "мышь": "mouse",
  "пульт": "remote",
  "клавиатура": "keyboard",
  "мобильный телефон": "cell phone",
  "микроволновка": "microwave",
  "духовка": "oven",
  "тостер": "toaster",
  "раковина": "sink",
  "холодильник": "refrigerator",
  "книга": "book",
  "часы": "clock",
  "ваза": "vase",
  "ножницы": "scissors",
  "плюшевый мишка": "teddy bear",
};

const YOLO_OBJECTS = [
  "лицо", "человек", "велосипед", "автомобиль", "мотоцикл", "самолет", "автобус", "поезд", "грузовик", "лодка",
  "светофор", "пожарный гидрант", "знак стоп", "парковочный счетчик", "скамейка", "птица", "кот", "собака",
  "лошадь", "овца", "корова", "слон", "медведь", "зебра", "жираф", "рюкзак", "зонт", "сумка", "галстук",
  "чемодан", "фрисби", "лыжи", "сноуборд", "спортивный мяч", "воздушный змей", "бейсбольная бита",
  "бейсбольная перчатка", "скейтборд", "доска для серфинга", "теннисная ракетка", "бутылка", "бокал для вина",
  "чашка", "вилка", "нож", "ложка", "миска", "банан", "яблоко", "бутерброд", "апельсин", "брокколи", "морковь",
  "хот-дог", "пицца", "пончик", "торт", "стул", "диван", "горшечное растение", "кровать", "обеденный стол",
  "туалет", "телевизор", "ноутбук", "мышь", "пульт", "клавиатура", "мобильный телефон", "микроволновка",
  "духовка", "тостер", "раковина", "холодильник", "книга", "часы", "ваза", "ножницы", "плюшевый мишка",
];

export default function ProcessPage() {
  const [dragActive, setDragActive] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [blurIntensity, setBlurIntensity] = useState([50]);
  const [blurType, setBlurType] = useState("blur");
  const [processing, setProcessing] = useState(false);
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [showRegistrationModal, setShowRegistrationModal] = useState(false);
  const [selectedObjects, setSelectedObjects] = useState<string[]>(["человек", "автомобиль"]);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("all");
  const [files, setFiles] = useState<any[]>([]);
  const [rateRemaining, setRateRemaining] = useState<number | null>(null);
  const [currentFileId, setCurrentFileId] = useState<string | null>(null);
  const [fileStatus, setFileStatus] = useState<string>("");
  const { token, isAuthenticated, user, logout } = useAuth();

  const objectCategories = {
    people: ["лицо", "человек"],
    vehicles: ["велосипед", "автомобиль", "мотоцикл", "самолет", "автобус", "поезд", "грузовик", "лодка"],
    animals: ["птица", "кот", "собака", "лошадь", "овца", "корова", "слон", "медведь", "зебра", "жираф"],
    objects: [
      "рюкзак", "зонт", "сумка", "галстук", "чемодан", "бутылка", "бокал для вина", "чашка", "книга",
      "мобильный телефон", "ноутбук", "пульт", "клавиатура", "мышь",
    ],
    furniture: ["стул", "диван", "горшечное растение", "кровать", "обеденный стол", "туалет"],
    food: ["банан", "яблоко", "бутерброд", "апельсин", "брокколи", "морковь", "хот-дог", "пицца", "пончик", "торт"],
    sports: [
      "фрисби", "лыжи", "сноуборд", "спортивный мяч", "воздушный змей", "бейсбольная бита",
      "бейсбольная перчатка", "скейтборд", "доска для серфинга", "теннисная ракетка",
    ],
    signs: ["светофор", "пожарный гидрант", "знак стоп", "парковочный счетчик"],
    appliances: ["телевизор", "микроволновка", "духовка", "тостер", "раковина", "холодильник"],
    other: ["скамейка", "вилка", "нож", "ложка", "миска", "часы", "ваза", "ножницы", "плюшевый мишка"],
  };

  const handleSwitchToLogin = () => {
    setShowRegistrationModal(false);
    setTimeout(() => setShowLoginModal(true), 100);
  };

  const handleSwitchToRegister = () => {
    setShowLoginModal(false);
    setTimeout(() => setShowRegistrationModal(true), 100);
  };

  const toggleObject = (object: string) => {
    setSelectedObjects((prev) => (prev.includes(object) ? prev.filter((o) => o !== object) : [...prev, object]));
  };

  const getFilteredObjects = () => {
    const objects =
      selectedCategory === "all"
        ? YOLO_OBJECTS
        : objectCategories[selectedCategory as keyof typeof objectCategories] || [];
    return objects.filter((obj) => obj.toLowerCase().includes(searchTerm.toLowerCase()));
  };

  const selectAllInCategory = () => {
    const categoryObjects =
      selectedCategory === "all"
        ? YOLO_OBJECTS
        : objectCategories[selectedCategory as keyof typeof objectCategories] || [];
    setSelectedObjects((prev) => [...new Set([...prev, ...categoryObjects])]);
  };

  const clearAllInCategory = () => {
    const categoryObjects =
      selectedCategory === "all"
        ? YOLO_OBJECTS
        : objectCategories[selectedCategory as keyof typeof objectCategories] || [];
    setSelectedObjects((prev) => prev.filter((obj) => !categoryObjects.includes(obj)));
  };

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0];
      if (file.size > 50 * 1024 * 1024) {
        alert("Файл превышает лимит в 50 MB");
        return;
      }
      setUploadedFile(file);
    }
  }, []);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.size > 50 * 1024 * 1024) {
        alert("Файл превышает лимит в 50 MB");
        return;
      }
      setUploadedFile(file);
    }
  };

const handleProcess = async () => {
  if (!uploadedFile || selectedObjects.length === 0) {
    alert("Выберите файл и хотя бы один объект для обработки");
    return;
  }

  setProcessing(true);
  setFileStatus("⏳ Загрузка...");

  const formData = new FormData();
  formData.append("file", uploadedFile);
  formData.append("blur_type", blurType === "blur" ? "gaussian" : blurType);
  formData.append("intensity", Math.round(blurIntensity[0] / 10).toString());
  formData.append("object_types", selectedObjects.map((obj) => RUS_TO_ENG_MAPPING[obj] || obj).join(","));

  try {
    const res = await fetch(`${API_BASE}/api/upload`, {
      method: "POST",
      body: formData,
      headers: isAuthenticated && token ? { Authorization: `Bearer ${token}` } : undefined,
    });

    if (!res.ok) {
      const error = await res.json().catch(() => ({}));
      if (res.status === 401 && isAuthenticated) {
        alert("Сессия истекла, пожалуйста, войдите снова");
        logout();
      } else if (res.status === 429 && !isAuthenticated) {
        alert("Исчерпан лимит в 3 попытки. Попробуйте завтра или войдите в аккаунт.");
        setFileStatus("❌ Лимит попыток исчерпан");
        return;
      }
      throw new Error(error.message || `Ошибка загрузки: ${res.status}`);
    }

    const data = await res.json();
    console.log("Upload response:", data);
    const fileId = data.data?.id;
    if (!fileId) {
      throw new Error("Не получен ID файла от сервера");
    }

    const remaining = !isAuthenticated ? res.headers.get("X-RateLimit-Remaining") : null;
    if (remaining !== null) {
      setRateRemaining(parseInt(remaining));
    }

    if (isAuthenticated) {
      setFileStatus("⏳ Обрабатывается...");
      pollStatus(fileId);
      fetchFiles();
    } else {
      // Эмуляция прогресса для анонимных пользователей
      setFileStatus("⏳ Обрабатывается...");
      let progress = 0;
      const progressInterval = setInterval(() => {
        progress += 20;
        setFileStatus(`⏳ Обрабатывается... ${progress}%`);
        if (progress >= 100) {
          clearInterval(progressInterval);
          setFileStatus("✅ Обработка завершена!");
          setCurrentFileId(fileId);
        }
      }, 900); // ~4.5 секунды
    }
  } catch (error: any) {
    console.error("Ошибка загрузки:", error);
    setFileStatus(`❌ Ошибка: ${error.message || "Не удалось загрузить файл"}`);
  } finally {
    setProcessing(false);
  }
};

const pollStatus = async (fileId: string) => {
  try {
    const res = await fetch(`${API_BASE}/api/files/${fileId}`, {
      method: "GET",
      headers: isAuthenticated && token ? { Authorization: `Bearer ${token}` } : undefined,
    });

    if (!res.ok) {
      const error = await res.json().catch(() => ({}));
      console.error("Polling error response:", error, "Status:", res.status); // Логируем ошибку
      if (res.status === 401 && isAuthenticated) {
        alert("Сессия истекла, пожалуйста, войдите снова");
        logout();
      } else if (res.status === 403 && !isAuthenticated) {
        setFileStatus("❌ Ошибка: Сервер отклонил запрос на проверку статуса. Попробуйте позже.");
        setCurrentFileId(null);
        return;
      }
      throw new Error(error.message || `Ошибка получения статуса: ${res.status}`);
    }

    const { data } = await res.json();
    console.log("Polling response:", data); // Логируем ответ
    if (!data?.status) {
      throw new Error("Статус файла не получен");
    }

    setFileStatus(data.status);
    if (data.status === "completed") {
      setFileStatus("✅ Обработка завершена!");
      setCurrentFileId(fileId);
    } else if (data.status === "failed") {
      setFileStatus(`❌ Ошибка: ${data.error_message || "Неизвестная ошибка"}`);
      setCurrentFileId(null);
    } else {
      setTimeout(() => pollStatus(fileId), 2500);
    }
  } catch (error: any) {
    console.error("Ошибка polling:", error);
    setFileStatus(`❌ Ошибка: ${error.message || "Не удалось проверить статус"}`);
    setCurrentFileId(null);
  }
};

  const fetchFiles = async () => {
    if (!isAuthenticated || !token) return;
    try {
      const res = await fetch(`${API_BASE}/api/files`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data = await res.json();
        setFiles(data.data);
      }
    } catch (error) {
      console.error("Error fetching files:", error);
    }
  };

  const handleDownload = async (fileId: string, type: "original" | "processed" = "processed") => {
    try {
      const res = await fetch(`${API_BASE}/api/files/${fileId}?type=${type}`, {
        method: "GET",
        headers: isAuthenticated && token ? { Authorization: `Bearer ${token}` } : undefined,
      });

      if (!res.ok) {
        const error = await res.json().catch(() => ({}));
        if (res.status === 401 && isAuthenticated) {
          alert("Сессия истекла, пожалуйста, войдите снова");
          logout();
        }
        throw new Error(error.message || `Ошибка скачивания: ${res.status}`);
      }

      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `${type}-${fileId}.${uploadedFile?.name.split(".").pop() || "file"}`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (error: any) {
      console.error(`Ошибка скачивания (${type}):`, error);
      alert(`Ошибка: ${error.message || "Не удалось скачать файл"}`);
    }
  };

  const handleDelete = async (fileId: string) => {
    if (!isAuthenticated) return;
    try {
      const res = await fetch(`${API_BASE}/api/files/${fileId}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        fetchFiles();
      } else {
        alert("Ошибка удаления");
      }
    } catch (error) {
      console.error("Delete error:", error);
      alert("Ошибка соединения");
    }
  };

  useEffect(() => {
    if (isAuthenticated) {
      fetchFiles();
    }
    const createParticles = () => {
      const particlesContainer = document.getElementById("particles");
      if (!particlesContainer) return;

      for (let i = 0; i < 30; i++) {
        const particle = document.createElement("div");
        particle.className = "particle";
        particle.style.left = Math.random() * 100 + "%";
        particle.style.animationDelay = Math.random() * 15 + "s";
        particle.style.animationDuration = Math.random() * 10 + 10 + "s";
        particlesContainer.appendChild(particle);
      }
    };

    createParticles();
  }, [isAuthenticated]);

  return (
    <div className="min-h-screen bg-black relative">
      <div className="bg-animation absolute top-0 left-0 w-full h-full z-0"></div>
      <div className="particles absolute top-0 left-0 w-full h-full z-10" id="particles"></div>

      <header className="fixed top-0 w-full z-50 backdrop-blur-lg bg-black/10 border-b border-white/10 transition-all duration-300">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-4">
              <Link href="/" className="flex items-center gap-2 text-white/80 hover:text-white transition-colors">
                <ArrowLeft className="w-5 h-5" />
              </Link>
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
                    <Button className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 shadow-lg hover:shadow-purple-500/25 transition-all duration-300 transform hover:-translate-y-1">
                      Регистрация
                    </Button>
                  </RegistrationModal>
                </>
              )}
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-24 pb-12 relative z-20">
        <h1 className="font-geist font-bold text-3xl lg:text-4xl text-white mb-8">Обработать файл</h1>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* Upload Section */}
          <Card className="bg-white/5 backdrop-blur-sm border-white/10">
            <CardHeader>
              <CardTitle className="font-geist flex items-center gap-2 text-white text-xl">
                <Upload className="w-6 h-6" />
                Загрузка файла
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div
                className={`border-2 border-dashed rounded-lg p-8 text-center ${
                  dragActive ? "border-purple-400 bg-purple-400/10" : "border-white/20"
                }`}
                onDragEnter={handleDrag}
                onDragOver={handleDrag}
                onDragLeave={handleDrag}
                onDrop={handleDrop}
              >
                <input
                  type="file"
                  id="file-upload"
                  className="hidden"
                  accept="image/jpeg,image/png,image/gif,image/webp,image/bmp,image/tiff,video/mp4,video/avi,video/mov,video/webm,video/mkv,video/wmv,video/flv"
                  onChange={handleFileSelect}
                />
                <label htmlFor="file-upload" className="cursor-pointer">
                  <div className="flex flex-col items-center gap-4">
                    <Upload className="w-12 h-12 text-purple-400" />
                    <p className="font-manrope text-white">
                      {uploadedFile ? uploadedFile.name : "Перетащите файл или нажмите для выбора"}
                    </p>
                    <p className="font-manrope text-xs text-white/60">
                      Поддерживаемые форматы: JPG, PNG, GIF, WebP, BMP, TIFF, MP4, AVI, MOV, WebM, MKV, WMV, FLV
                    </p>
                    <p className="font-manrope text-xs text-white/60">Максимум 50 MB</p>
                  </div>
                </label>
              </div>
              {rateRemaining !== null && !isAuthenticated && (
                <p className="font-manrope text-sm text-white/60 mt-4">
                  Осталось загрузок сегодня: {rateRemaining}
                </p>
              )}
            </CardContent>
          </Card>

          {/* Object Selection */}
          <Card className="bg-white/5 backdrop-blur-sm border-white/10">
            <CardHeader>
              <CardTitle className="font-geist flex items-center gap-2 text-white text-xl">
                <Search className="w-6 h-6" />
                Выберите объекты
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-3">
                <Select value={selectedCategory} onValueChange={setSelectedCategory}>
                  <SelectTrigger className="font-manrope bg-white/10 border-white/20 text-white">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-black/80 backdrop-blur-lg border-white/20 shadow-2xl">
                    <SelectItem value="all" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Все
                    </SelectItem>
                    <SelectItem value="people" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Люди
                    </SelectItem>
                    <SelectItem value="vehicles" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Транспорт
                    </SelectItem>
                    <SelectItem value="animals" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Животные
                    </SelectItem>
                    <SelectItem value="objects" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Объекты
                    </SelectItem>
                    <SelectItem value="furniture" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Мебель
                    </SelectItem>
                    <SelectItem value="food" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Еда
                    </SelectItem>
                    <SelectItem value="sports" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Спорт
                    </SelectItem>
                    <SelectItem value="signs" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Знаки
                    </SelectItem>
                    <SelectItem value="appliances" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Бытовая техника
                    </SelectItem>
                    <SelectItem value="other" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                      Прочее
                    </SelectItem>
                  </SelectContent>
                </Select>
                <Input
                  placeholder="Поиск объектов..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="font-manrope bg-white/10 border-white/20 text-white placeholder:text-white/50"
                />
              </div>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={selectAllInCategory}
                  className="font-manrope bg-white/10 border-white/20 text-white hover:bg-white/20"
                >
                  Выбрать все
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={clearAllInCategory}
                  className="font-manrope bg-white/10 border-white/20 text-white hover:bg-white/20"
                >
                  Очистить
                </Button>
              </div>
              <div className="max-h-80 overflow-y-auto space-y-2 pr-2">
                <div className="grid grid-cols-2 gap-2">
                  {getFilteredObjects().map((object) => (
                    <div
                      key={object}
                      onClick={() => toggleObject(object)}
                      className={`group relative p-3 rounded-lg cursor-pointer transition-all duration-200 border ${
                        selectedObjects.includes(object)
                          ? "bg-gradient-to-br from-purple-500/30 to-blue-500/30 border-purple-400/50 shadow-lg shadow-purple-500/20"
                          : "bg-white/5 hover:bg-white/10 border-white/10 hover:border-white/20"
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <span className="font-manrope text-sm text-white capitalize leading-tight">{object}</span>
                        <div
                          className={`w-4 h-4 rounded-full border-2 flex items-center justify-center transition-all ${
                            selectedObjects.includes(object)
                              ? "border-purple-400 bg-purple-400"
                              : "border-white/30 group-hover:border-white/50"
                          }`}
                        >
                          {selectedObjects.includes(object) && <Check className="w-2.5 h-2.5 text-white" />}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
              <div className="pt-2 border-t border-white/10">
                <div className="flex items-center justify-between">
                  <p className="font-manrope text-xs text-white/60">Выбрано: {selectedObjects.length} объектов</p>
                  {selectedObjects.length > 0 && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setSelectedObjects([])}
                      className="font-manrope text-xs text-white/60 hover:text-white hover:bg-white/10 p-1"
                    >
                      Очистить все
                    </Button>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Settings and Export */}
          <div className="space-y-6">
            <Card className="bg-white/5 backdrop-blur-sm border-white/10">
              <CardHeader>
                <CardTitle className="font-geist flex items-center gap-2 text-white text-xl">
                  <Settings className="w-6 h-6" />
                  Настройки размытия
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div>
                  <Label className="font-manrope text-sm font-medium mb-3 block text-white">Тип эффекта</Label>
                  <Select value={blurType} onValueChange={setBlurType}>
                    <SelectTrigger className="font-manrope bg-white/10 border-white/20 text-white">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent className="bg-black/80 backdrop-blur-lg border-white/20 shadow-2xl">
                      <SelectItem value="blur" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                        Размытие
                      </SelectItem>
                      <SelectItem value="pixelate" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                        Пикселизация
                      </SelectItem>
                      <SelectItem value="mask" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                        Цветная маска
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label className="font-manrope text-sm font-medium mb-3 block text-white">
                    Интенсивность: {blurIntensity[0]}%
                  </Label>
                  <Slider
                    value={blurIntensity}
                    onValueChange={setBlurIntensity}
                    max={100}
                    min={10}
                    step={5}
                    className="w-full"
                  />
                </div>
                <div className="space-y-3">
                  <Label className="font-manrope text-sm font-medium text-white">Выбранные объекты</Label>
                  <div className="flex flex-wrap gap-2 max-h-24 overflow-y-auto">
                    {selectedObjects.map((object) => (
                      <Badge
                        key={object}
                        variant="secondary"
                        className="font-manrope bg-gradient-to-r from-purple-500/20 to-blue-500/20 text-white border border-purple-400/30"
                      >
                        {object}
                      </Badge>
                    ))}
                  </div>
                </div>
                <Button
                  onClick={handleProcess}
                  disabled={!uploadedFile || processing || selectedObjects.length === 0}
                  className="w-full font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                  size="lg"
                >
                  {processing ? "Обработка..." : "Применить"}
                </Button>
                {processing && (
                  <div className="w-full bg-white/10 rounded-full h-2">
                    <div
                      className="bg-gradient-to-r from-purple-500 to-blue-500 h-2 rounded-full animate-pulse"
                      style={{ width: "60%" }}
                    ></div>
                  </div>
                )}
              </CardContent>
            </Card>
            <Card className="bg-white/5 backdrop-blur-sm border-white/10">
              <CardHeader>
                <CardTitle className="font-geist flex items-center gap-2 text-white text-xl">
                  <Download className="w-6 h-6" />
                  Экспорт
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <Label className="font-manrope text-sm font-medium mb-2 block text-white">Качество</Label>
                    <Select defaultValue="original">
                      <SelectTrigger className="font-manrope bg-white/10 border-white/20 text-white">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent className="bg-black/80 backdrop-blur-lg border-white/20 shadow-2xl">
                        <SelectItem value="original" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                          Оригинальное
                        </SelectItem>
                        <SelectItem value="high" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                          Высокое
                        </SelectItem>
                        <SelectItem value="medium" className="font-manrope text-white hover:bg-white/10 focus:bg-white/10">
                          Среднее
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <Button
                    onClick={() => currentFileId && handleDownload(currentFileId, "processed")}
                    disabled={!currentFileId || processing}
                    className="w-full font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                  >
                    Скачать результат
                  </Button>
                  {fileStatus && <p className="font-manrope text-sm text-white/60">{fileStatus}</p>}
                </div>
              </CardContent>
            </Card>
            {isAuthenticated && files.length > 0 && (
              <Card className="bg-white/5 backdrop-blur-sm border-white/10">
                <CardHeader>
                  <CardTitle className="font-geist text-white text-xl">Ваши файлы</CardTitle>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-2">
                    {files.map((file) => (
                      <li key={file.id} className="flex justify-between items-center">
                        <div className="flex flex-col">
                          <span className="font-manrope text-white text-sm">{file.original_name}</span>
                          <span className="font-manrope text-xs text-white/60">
                            {file.status === "completed" && "✅ Завершено"}
                            {file.status === "processing" && "⏳ Обрабатывается"}
                            {file.status === "failed" && `❌ Ошибка: ${file.error_message || "Неизвестная ошибка"}`}
                          </span>
                        </div>
                        <div className="flex gap-2">
                          {file.status === "completed" && (
                            <>
                              <Button
                                size="sm"
                                onClick={() => handleDownload(file.id, "original")}
                                className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                              >
                                Оригинал
                              </Button>
                              <Button
                                size="sm"
                                onClick={() => handleDownload(file.id, "processed")}
                                className="font-manrope bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600"
                              >
                                Обработанный
                              </Button>
                            </>
                          )}
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => handleDelete(file.id)}
                            className="font-manrope"
                          >
                            <Trash className="w-4 h-4" />
                          </Button>
                        </div>
                      </li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}