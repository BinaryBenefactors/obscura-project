"use client";

import type React from "react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Check, Camera } from "lucide-react";
import { useAuth } from "@/components/AuthContext";

const API_LINK = process.env.NEXT_PUBLIC_API_LINK || "http://localhost:8080";

interface RegistrationModalProps {
  children: React.ReactNode;
  onSwitchToLogin?: () => void;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function RegistrationModal({
  children,
  onSwitchToLogin,
  open: externalOpen,
  onOpenChange,
}: RegistrationModalProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [internalOpen, setInternalOpen] = useState(false);
  const { login } = useAuth();

  const open = externalOpen !== undefined ? externalOpen : internalOpen;
  const setOpen = onOpenChange || setInternalOpen;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (name.length < 2) {
      alert("Имя должно быть минимум 2 символа");
      return;
    }
    if (password.length < 6) {
      alert("Пароль должен быть минимум 6 символов");
      return;
    }
    try {
      const res = await fetch(API_LINK + "/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, email, password }),
      });
      if (res.ok) {
        const data = await res.json();
        login(data.token, { name, email });
        setOpen(false);
      } else {
        const error = await res.json();
        alert(error.message || "Ошибка регистрации");
      }
    } catch (error) {
      console.error("Registration error:", error);
      alert("Ошибка соединения");
    }
  };

  const handleSwitchToLogin = () => {
    setOpen(false);
    if (onSwitchToLogin) {
      setTimeout(() => onSwitchToLogin(), 100);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="modal-content max-h-[90vh] z-1001 overflow-y-auto p-4 sm:p-6 text-gray-900">
        <DialogHeader>
          <div className="flex items-center gap-2 mb-4">
            <div className="w-8 h-8 bg-gradient-to-r from-cyan-400 to-blue-500 rounded-lg flex items-center justify-center">
              <Camera className="w-5 h-5 text-white" />
            </div>
            <span className="font-geist font-semibold text-xl text-gray-900">Obscura</span>
          </div>
          <DialogTitle className="font-geist text-2xl text-gray-900">
            Доступ к профессиональным инструментам
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-8">
          <div className="space-y-4">
            <p className="font-manrope text-gray-600">Зарегистрируйтесь, чтобы:</p>
            <ul className="space-y-3">
              <li className="flex items-center gap-2 font-manrope text-sm text-gray-700">
                <Check className="w-4 h-4 text-cyan-500" />
                Сохранять историю обработок
              </li>
              <li className="flex items-center gap-2 font-manrope text-sm text-gray-700">
                <Check className="w-4 h-4 text-cyan-500" />
                Получать приоритетную поддержку
              </li>
              <li className="flex items-center gap-2 font-manrope text-sm text-gray-700">
                <Check className="w-4 h-4 text-cyan-500" />
                Экспортировать файлы без водяных знаков
              </li>
            </ul>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name" className="font-manrope text-gray-700">
                Имя
              </Label>
              <Input
                id="name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="font-manrope bg-white border-gray-300 text-gray-900 placeholder:text-gray-400 focus:border-cyan-400 py-3"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email" className="font-manrope text-gray-700">
                Email
              </Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="font-manrope bg-white border-gray-300 text-gray-900 placeholder:text-gray-400 focus:border-cyan-400 py-3"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password" className="font-manrope text-gray-700">
                Пароль
              </Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="font-manrope bg-white border-gray-300 text-gray-900 placeholder:text-gray-400 focus:border-cyan-400 py-3"
                required
              />
            </div>
            <Button
              type="submit"
              className="w-full font-manrope bg-gradient-to-r from-cyan-400 to-blue-500 hover:from-cyan-500 hover:to-blue-600 text-white border-0 py-3"
            >
              Создать аккаунт (это бесплатно)
            </Button>
          </form>

          <p className="text-xs text-gray-500 text-center font-manrope">
            Уже есть аккаунт?{" "}
            <button
              type="button"
              onClick={handleSwitchToLogin}
              className="text-cyan-500 hover:text-cyan-600 hover:underline"
            >
              Войти
            </button>
            <br />
            Регистрируясь, вы соглашаетесь с{" "}
            <a href="#" className="text-cyan-500 hover:text-cyan-600 hover:underline">
              условиями использования
            </a>{" "}
            и{" "}
            <a href="#" className="text-cyan-500 hover:text-cyan-600 hover:underline">
              политикой конфиденциальности
            </a>
          </p>
        </div>
      </DialogContent>
    </Dialog>
  );
}