import type React from "react";
import type { Metadata } from "next";
import { Geist, Manrope } from "next/font/google";
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
    <html lang="ru" className={`${geist.variable} ${manrope.variable} antialiased`}>
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}