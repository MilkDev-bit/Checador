/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      opacity: {
        '2': '0.02',
        '3': '0.03',
        '4': '0.04',
        '8': '0.08',
        '12': '0.12',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif']
      },
      colors: {
        brand: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554',
        },
        surface: {
          DEFAULT: '#09090b', // Zinc-950
          50: '#fafafa',
          100: '#f4f4f5',
          200: '#e4e4e7',
          300: '#d4d4d8',
          400: '#a1a1aa',
          500: '#71717a',
          600: '#52525b',
          700: '#3f3f46',
          800: '#27272a',
          900: '#18181b',
          950: '#09090b',
        }
      },
      backgroundImage: {
        'gradient-brand': 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 50%, #0ea5e9 100%)',
        'gradient-card': 'linear-gradient(145deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0.01) 100%)',
        'gradient-card-light': 'linear-gradient(145deg, rgba(255,255,255,0.9) 0%, rgba(255,255,255,0.4) 100%)',
        'gradient-success': 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
        'gradient-danger': 'linear-gradient(135deg, #ef4444 0%, #e11d48 100%)',
        'gradient-dark': 'linear-gradient(135deg, #18181b 0%, #09090b 100%)',
        'primary-mesh': 'radial-gradient(at 40% 20%, hsla(228,100%,74%,0.15) 0px, transparent 50%), radial-gradient(at 80% 0%, hsla(189,100%,56%,0.15) 0px, transparent 50%), radial-gradient(at 0% 50%, hsla(355,100%,93%,0.1) 0px, transparent 50%)',
      },
      animation: {
        'fade-in': 'fadeIn 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
        'slide-up': 'slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1)',
        'slide-down': 'slideDown 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'float': 'float 6s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: { from: { opacity: '0' }, to: { opacity: '1' } },
        slideUp: { from: { opacity: '0', transform: 'translateY(12px)' }, to: { opacity: '1', transform: 'translateY(0)' } },
        slideDown: { from: { opacity: '0', transform: 'translateY(-12px)' }, to: { opacity: '1', transform: 'translateY(0)' } },
        float: { '0%,100%': { transform: 'translateY(0)' }, '50%': { transform: 'translateY(-6px)' } },
      },
      boxShadow: {
        'glass': '0 8px 32px -4px rgba(0, 0, 0, 0.1)',
        'glass-dark': '0 8px 32px -4px rgba(0, 0, 0, 0.3)',
        'float': '0 14px 28px rgba(0,0,0,0.05), 0 10px 10px rgba(0,0,0,0.02)',
        'float-dark': '0 14px 28px rgba(0,0,0,0.3), 0 10px 10px rgba(0,0,0,0.2)',
        'glow': '0 0 20px rgba(59, 130, 246, 0.3)',
        'glow-brand': '0 0 24px -4px rgba(99, 102, 241, 0.4)',
        'inner-light': 'inset 0 1px 0 0 rgba(255, 255, 255, 0.15)',
        'inner-border': 'inset 0 0 0 1px rgba(0, 0, 0, 0.05)',
        'inner-border-dark': 'inset 0 0 0 1px rgba(255, 255, 255, 0.05)',
      }
    }
  },
  plugins: []
}
