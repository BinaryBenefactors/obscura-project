"use client"

import * as React from "react"
import * as SliderPrimitive from "@radix-ui/react-slider"
import { cn } from "@/lib/utils"

interface SliderProps extends React.ComponentProps<typeof SliderPrimitive.Root> {
  min?: number
  max?: number
}

function Slider({
  className,
  defaultValue,
  value,
  min = 0,
  max = 100,
  onValueChange,
  ...props
}: SliderProps) {
  const _values = React.useMemo(
    () =>
      Array.isArray(value)
        ? value
        : Array.isArray(defaultValue)
        ? defaultValue
        : [min, max],
    [value, defaultValue, min, max]
  )

  // Обновление кастомного курсора при перемещении Thumb
  const handleValueChange = (val: number[]) => {
    onValueChange?.(val)

    const cursor = document.querySelector(".cursor") as HTMLElement
    const cursorFollower = document.querySelector(".cursor-follower") as HTMLElement
    const sliderThumb = document.querySelector('[data-slot="slider-thumb"]') as HTMLElement

    if (sliderThumb && cursor && cursorFollower) {
      const rect = sliderThumb.getBoundingClientRect()
      const x = rect.left + rect.width / 2
      const y = rect.top + rect.height / 2

      cursor.style.transform = `translate(${x}px, ${y}px)`
      cursorFollower.style.transform = `translate(${x}px, ${y}px)`
    }
  }

  // Обработка клика по треку
  const handleTrackPointerDown = (e: React.PointerEvent<HTMLDivElement>) => {
    const rect = e.currentTarget.getBoundingClientRect()
    let newValue: number
    if (rect.width >= rect.height) {
      // горизонтальный слайдер
      const clickPosition = e.clientX - rect.left
      newValue = min + (clickPosition / rect.width) * (max - min)
    } else {
      // вертикальный слайдер
      const clickPosition = rect.bottom - e.clientY
      newValue = min + (clickPosition / rect.height) * (max - min)
    }
    onValueChange?.([Math.round(newValue)])
  }

  return (
    <SliderPrimitive.Root
      data-slot="slider"
      defaultValue={defaultValue}
      value={value}
      min={min}
      max={max}
      onValueChange={handleValueChange}
      className={cn(
        "relative flex w-full touch-none items-center select-none data-[disabled]:opacity-50 data-[orientation=vertical]:h-72 data-[orientation=vertical]:w-auto data-[orientation=vertical]:flex-col",
        className
      )}
      {...props}
    >
      <SliderPrimitive.Track
        data-slot="slider-track"
        className={cn(
          "bg-muted relative grow overflow-hidden rounded-full data-[orientation=horizontal]:h-3 data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-72 data-[orientation=vertical]:w-2"
        )}
        onPointerDown={handleTrackPointerDown}
      >
        <SliderPrimitive.Range
          data-slot="slider-range"
          className={cn(
            "bg-primary absolute data-[orientation=horizontal]:h-full data-[orientation=vertical]:w-full"
          )}
        />
      </SliderPrimitive.Track>

      {Array.from({ length: _values.length }, (_, index) => (
        <SliderPrimitive.Thumb
          data-slot="slider-thumb"
          key={index}
          className="border-primary bg-background ring-ring/50 block w-5 h-5 shrink-0 rounded-full border shadow-sm transition-[color,box-shadow] hover:ring-4 focus-visible:ring-4 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
        />
      ))}
    </SliderPrimitive.Root>
  )
}

export { Slider }
