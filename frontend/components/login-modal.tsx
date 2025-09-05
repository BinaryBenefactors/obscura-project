"use client";

import type React from "react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Camera } from "lucide-react";
import { useAuth } from "@/components/AuthContext";

const API_LINK = process.env.NEXT_PUBLIC_API_LINK || "http://localhost:8080";

interface LoginModalProps {
  children: React.ReactNode;
  onSwitchToRegister?: () => void;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function LoginModal({
  children,
  onSwitchToRegister,
  open: externalOpen,
  onOpenChange,
}: LoginModalProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [internalOpen, setInternalOpen] = useState(false);
  const { login } = useAuth();

  const open = externalOpen !== undefined ? externalOpen : internalOpen;
  const setOpen = onOpenChange || setInternalOpen;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password.length < 6) {
      alert("Пароль должен быть минимум 6 символов");
      return;
    }
    try {
      const res = await fetch(API_LINK + "/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
      if (res.ok) {
        const data = await res.json();
        const token = data.token;
        const name = data.user?.name || ""; // Извлекаем имя из ответа
        if (!token) throw new Error("Токен не получен");
        login(token, { email, name }); // Передаём email и name
        setOpen(false);
      } else {
        const error = await res.json();
        alert(error.message || "Ошибка входа");
      }
    } catch (error) {
      console.error("Login error:", error);
      alert("Ошибка соединения");
    }
  };

  const handleSwitchToRegister = () => {
    setOpen(false);
    if (onSwitchToRegister) {
      setTimeout(() => onSwitchToRegister(), 100);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-md bg-white/95 backdrop-blur-lg border border-gray-200 text-gray-900 p-8">
        <DialogHeader>
          <div className="flex items-center gap-2 mb-4">
            <div className="w-8 h-8 bg-gradient-to-r from-cyan-400 to-blue-500 rounded-lg flex items-center justify-center">
              <Camera className="w-5 h-5 text-white" />
            </div>
            <span className="font-geist font-semibold text-xl text-gray-900">Obscura</span>
          </div>
          <DialogTitle className="font-geist text-2xl text-gray-900">Вход в аккаунт</DialogTitle>
        </DialogHeader>

        <div className="space-y-8">
          <p className="font-manrope text-gray-600">
            Войдите в свой аккаунт для доступа к профессиональным инструментам
          </p>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="login-email" className="font-manrope text-gray-700">
                Email
              </Label>
              <Input
                id="login-email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="font-manrope bg-white border-gray-300 text-gray-900 placeholder:text-gray-400 focus:border-cyan-400 py-3"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="login-password" className="font-manrope text-gray-700">
                Пароль
              </Label>
              <Input
                id="login-password"
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
              Войти
            </Button>
          </form>

          <div className="text-center">
            <a href="#" className="text-sm text-cyan-500 hover:text-cyan-600 hover:underline font-manrope">
              Забыли пароль?
            </a>
          </div>

          <p className="text-xs text-gray-500 text-center font-manrope">
            Нет аккаунта?{" "}
            <button
              type="button"
              onClick={handleSwitchToRegister}
              className="text-cyan-500 hover:text-cyan-600 hover:underline"
            >
              Зарегистрируйтесь бесплатно
            </button>
          </p>
        </div>
      </DialogContent>
    </Dialog>
  );
}