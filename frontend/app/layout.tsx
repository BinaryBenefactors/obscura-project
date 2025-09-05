import type React from "react";
import type { Metadata } from "next";
import { Geist, Manrope, Mitr } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/components/AuthContext";


const geist = Geist({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-geist",
});

const manrope = Manrope({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-manrope",
});

const mitr = Mitr({
  weight: ['200', '300', '400', '500', '600', '700'], // Укажите нужные веса
  subsets: ['latin'], // Или 'latin-ext', если нужно
  display: 'swap', // Оптимизация загрузки
  variable: '--font-mitr',
});

export const metadata: Metadata = {
  title: "AI Blur - Автоматическое размытие конфиденциальных объектов",
  description:
    "Сервис для фотографов и видеографов, который за секунды скрывает лица, номера машин или любые объекты в кадре с помощью AI.",
  generator: "v0.app",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ru" className={`${geist.variable} ${manrope.variable} ${mitr.variable} antialiased`}>
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}