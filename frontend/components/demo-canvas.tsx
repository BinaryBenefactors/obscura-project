import { useEffect, useRef, useState } from "react"

export function DemoCanvas() {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [blurIntensity, setBlurIntensity] = useState(10)
  const [effectType, setEffectType] = useState("blur") // blur | pixelate | motion
  const [hiddenObject, setHiddenObject] = useState("face") // face | laptop

  // demo image
  const imgUrl =
    "https://images.unsplash.com/photo-1573164713988-8665fc963095?w=800&h=500&fit=crop"

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext("2d")
    if (!ctx) return

    const img = new Image()
    img.crossOrigin = "anonymous"
    img.src = imgUrl
    img.onload = () => {
      canvas.width = img.width
      canvas.height = img.height

      // рисуем базовое изображение
      ctx.drawImage(img, 0, 0)

      // определим зоны для разных объектов
      const objects: Record<
        string,
        { x: number; y: number; w: number; h: number }[]
      > = {
        face: [
          { x: 545, y: 90, w: 50, h: 60 },  // первое лицо
          { x: 445, y: 100, w: 50, h: 60 },  // второе лицо (левее)
        ],
        laptop: [
          { x: 560, y: 170, w: 100, h: 60 },
          { x: 395, y: 185, w: 100, h: 55 },
        ],
      }

      const zones = objects[hiddenObject] || []

      zones.forEach((zone) => {
        if (effectType === "blur") {
          const offCanvas = document.createElement("canvas")
          const offCtx = offCanvas.getContext("2d")!
          offCanvas.width = img.width
          offCanvas.height = img.height

          const blur = Math.pow(blurIntensity / 100, 2) * 20
          offCtx.filter = `blur(${blur}px)`
          offCtx.drawImage(img, 0, 0)

          const blurredPart = offCtx.getImageData(zone.x, zone.y, zone.w, zone.h)
          ctx.putImageData(blurredPart, zone.x, zone.y)
        }

        if (effectType === "pixelate") {
          const pixels = blurIntensity / 5 + 2
          const offCanvas = document.createElement("canvas")
          const offCtx = offCanvas.getContext("2d")!
          offCanvas.width = zone.w
          offCanvas.height = zone.h
          offCtx.imageSmoothingEnabled = false
          offCtx.drawImage(
            img,
            zone.x,
            zone.y,
            zone.w,
            zone.h,
            0,
            0,
            zone.w / pixels,
            zone.h / pixels
          )
          ctx.imageSmoothingEnabled = false
          ctx.drawImage(
            offCanvas,
            0,
            0,
            zone.w / pixels,
            zone.h / pixels,
            zone.x,
            zone.y,
            zone.w,
            zone.h
          )
        }

        if (effectType === "motion") {
            const steps = Math.max(2, Math.floor(blurIntensity / 5)) // количество шагов
            const offCanvas = document.createElement("canvas")
            const offCtx = offCanvas.getContext("2d")!
            offCanvas.width = zone.w
            offCanvas.height = zone.h

            // копируем исходную область
            const objectData = ctx.getImageData(zone.x, zone.y, zone.w, zone.h)
            offCtx.putImageData(objectData, 0, 0)

            // смазываем горизонтально
            for (let i = 1; i <= steps; i++) {
                offCtx.globalAlpha = 1 / (i + 1)
                offCtx.drawImage(offCanvas, i, 0)
            }
            offCtx.globalAlpha = 1

            // вставляем размытую область обратно на основное изображение
            ctx.putImageData(offCtx.getImageData(0, 0, zone.w, zone.h), zone.x, zone.y)
            }

        // DEBUG рамка
        ctx.strokeStyle = "rgba(0, 150, 255, 0.8)"
        ctx.lineWidth = 3
        ctx.strokeRect(zone.x, zone.y, zone.w, zone.h)
      })
    }
  }, [blurIntensity, effectType, hiddenObject])

  return (
    <section className="demo" id="demo">
      <div className="container">
        <div className="section-header">
          <span className="section-tag">Демонстрация</span>
          <h2 className="section-title">Увидьте магию в действии</h2>
          <p className="section-subtitle">
            Настройте параметры и посмотрите, как работает Obscura
          </p>
        </div>

        <div className="demo-content flex flex-col md:flex-row gap-8">
          <div className="demo-visual flex-1 flex justify-center items-center">
            <canvas
              ref={canvasRef}
              className="border rounded-lg shadow-lg max-w-full h-auto"
            />
          </div>

          <div className="demo-controls flex-1">
            <div className="control-group mb-6">
              <label className="control-label">Интенсивность эффекта</label>
              <input
                type="range"
                min="0"
                max="100"
                value={blurIntensity}
                onChange={(e) => setBlurIntensity(Number(e.target.value))}
                className="slider w-full"
              />
            </div>

            <div className="control-group mb-6">
              <label className="control-label">Тип эффекта</label>
              <div className="toggle-group flex gap-2">
                <button
                  className={`toggle-btn ${effectType === "blur" ? "active" : ""}`}
                  onClick={() => setEffectType("blur")}
                >
                  Гауссово размытие
                </button>
                <button
                  className={`toggle-btn ${effectType === "motion" ? "active" : ""}`}
                  onClick={() => setEffectType("motion")}
                >
                  Движение
                </button>
                <button
                  className={`toggle-btn ${effectType === "pixelate" ? "active" : ""}`}
                  onClick={() => setEffectType("pixelate")}
                >
                  Пикселизация
                </button>
              </div>
            </div>

            <div className="control-group mb-6">
              <label className="control-label">Объекты для скрытия</label>
              <div className="toggle-group flex gap-2">
                <button
                  className={`toggle-btn ${hiddenObject === "face" ? "active" : ""}`}
                  onClick={() => setHiddenObject("face")}
                >
                  Лица
                </button>
                <button
                  className={`toggle-btn ${hiddenObject === "laptop" ? "active" : ""}`}
                  onClick={() => setHiddenObject("laptop")}
                >
                  Ноутбук
                </button>
              </div>
            </div>

            <a
              href="/process"
              className="btn-primary w-full inline-flex items-center justify-center gap-3 mt-8"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="white">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
              </svg>
              Применить эффекты
            </a>
          </div>
        </div>
      </div>
    </section>
  )
}
