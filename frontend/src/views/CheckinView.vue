<template>
  <div class="min-h-screen min-h-dvh bg-mesh flex flex-col md:flex-row">
    
    <!-- Sidebar / Top Navigation (Responsive) -->
    <aside class="md:w-64 glass-panel border-b md:border-b-0 md:border-r border-white/20 dark:border-white/5 flex-shrink-0 flex flex-col z-20 sticky top-0 md:h-screen md:sticky rounded-none">
      
      <!-- Brand & Header -->
      <div class="h-16 px-4 md:px-6 safe-top flex items-center justify-between border-b border-slate-200 dark:border-white/5" style="--header-base-padding: 0;">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0 bg-gradient-to-br from-brand-500 to-violet-600 shadow-glow-brand ring-2 ring-white/20 dark:ring-white/5">
            <CheckCircleIcon class="w-5 h-5 text-white" />
          </div>
          <div class="flex flex-col justify-center">
            <h1 class="font-bold text-base leading-tight tracking-tight text-slate-900 dark:text-white">PaseLista</h1>
            <p class="text-[11px] font-medium text-slate-500 dark:text-slate-400 truncate max-w-[120px]">{{ auth.user?.project_name || 'Mi Proyecto' }}</p>
          </div>
        </div>
        
        <!-- Mobile Actions -->
        <div class="flex items-center gap-2 md:hidden">
          <router-link to="/profile" class="w-8 h-8 rounded-full overflow-hidden flex-shrink-0 transition-transform active:scale-95 ring-2 ring-brand-100 dark:ring-brand-900/50">
            <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" class="w-full h-full object-cover" alt="Perfil" />
            <div v-else class="w-full h-full bg-brand-500 flex items-center justify-center text-white text-xs font-bold">
              {{ auth.user?.first_name?.[0] ?? '' }}{{ auth.user?.last_name?.[0] ?? '' }}
            </div>
          </router-link>
          <button v-if="!isPWA" @click="showMobileMenu = true" class="w-8 h-8 ml-2 flex items-center justify-center text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-surface-800 rounded-lg transition-colors">
            <Bars3BottomRightIcon class="w-6 h-6" />
          </button>
        </div>
      </div>

      <!-- Navigation Links (Hidden on small screens, shown on md up) -->
      <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-1.5 hidden md:block">
        <p class="px-3 text-[11px] font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-3">Menú Principal</p>
        
        <router-link to="/checkin" class="flex items-center gap-3 px-3 py-2.5 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400 font-semibold group transition-all">
          <MapPinIcon class="w-5 h-5 group-hover:scale-110 transition-transform" />
          Mi Check-in
        </router-link>
        
        <router-link to="/history" class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 hover:text-slate-900 dark:hover:text-white group transition-all">
          <ClockIcon class="w-5 h-5 group-hover:scale-110 transition-transform" />
          Historial
        </router-link>

        <router-link to="/qr" class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 hover:text-slate-900 dark:hover:text-white group transition-all">
          <QrCodeIcon class="w-5 h-5 group-hover:scale-110 transition-transform" />
          Escanear QR
        </router-link>
      </nav>

      <!-- Desktop Bottom Actions -->
      <div class="p-4 border-t border-slate-200 dark:border-white/5 hidden md:block space-y-3">
        <router-link to="/profile" class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
          <div class="w-8 h-8 rounded-full overflow-hidden flex-shrink-0 bg-brand-100 dark:bg-brand-900/40 border border-brand-200 dark:border-brand-800">
            <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" class="w-full h-full object-cover" alt="Perfil" />
            <div v-else class="w-full h-full text-brand-600 dark:text-brand-400 flex items-center justify-center text-xs font-bold">
              {{ auth.user?.first_name?.[0] ?? '' }}
            </div>
          </div>
          <div class="overflow-hidden">
            <p class="text-sm font-semibold text-slate-900 dark:text-white truncate">{{ auth.user?.first_name }} {{ auth.user?.last_name }}</p>
            <p class="text-xs text-slate-500 dark:text-slate-400 truncate">Ver perfil</p>
          </div>
        </router-link>
        
        <button @click="showLogoutModal = true" class="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-rose-500 dark:text-rose-400 font-semibold hover:bg-rose-50 dark:hover:bg-rose-500/10 transition-colors mt-2">
          <ArrowLeftOnRectangleIcon class="w-4 h-4" />
          Cerrar Sesión
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 relative overflow-y-auto flex flex-col md:pt-6 md:px-8 max-w-5xl mx-auto w-full">
      
      <!-- GPS background warning banner -->
      <Transition name="slide-banner">
        <div v-if="showHiddenWarning" class="glass-card mb-4 mx-4 md:mx-0 p-4 border border-amber-200 dark:border-amber-500/20 bg-amber-50/80 dark:bg-amber-900/20 flex items-start sm:items-center gap-4">
          <div class="p-2 bg-amber-100 dark:bg-amber-500/20 rounded-full text-amber-500 dark:text-amber-400">
            <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0" />
          </div>
          <div class="flex-1">
            <h4 class="text-amber-800 dark:text-amber-400 font-bold text-sm">GPS posiblemente pausado</h4>
            <p class="text-amber-700/80 dark:text-amber-300/70 text-xs mt-0.5 leading-relaxed">La app estuvo en segundo plano. Algunos puntos podrían no haberse registrado.</p>
          </div>
          <button @click="showHiddenWarning = false" class="p-1.5 hover:bg-amber-200/50 dark:hover:bg-amber-500/20 rounded-lg text-amber-600 dark:text-amber-400/60 transition-colors">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </Transition>

      <div class="px-4 py-6 md:py-0 space-y-6 flex-1">
        
        <!-- Header Section -->
        <header class="flex flex-col sm:flex-row sm:items-end justify-between gap-4 animate-slide-up">
          <div>
            <h2 class="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">Hola, {{ auth.user?.first_name }} 👋</h2>
            <p class="text-slate-500 dark:text-slate-400 font-medium mt-1">{{ todayStr }}</p>
          </div>
          <div class="glass-card px-5 py-3 border border-slate-200 dark:border-white/10 flex items-center justify-between gap-6 shadow-sm">
            <div class="flex flex-col">
              <span class="text-[10px] uppercase font-bold tracking-wider text-slate-400 dark:text-slate-500">Hora actual</span>
              <span class="text-xl font-bold font-mono text-brand-600 dark:text-brand-400">{{ currentTime }}</span>
            </div>
            <div class="h-10 w-px bg-slate-200 dark:bg-white/10"></div>
            <div class="flex flex-col items-end">
              <span class="text-[10px] uppercase font-bold tracking-wider text-slate-400 dark:text-slate-500">Status GPS</span>
              <div v-if="isTracking" class="flex items-center gap-1.5 mt-1">
                <span class="relative flex h-2.5 w-2.5">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
                </span>
                <span class="text-sm font-semibold text-emerald-600 dark:text-emerald-400">Activo</span>
              </div>
              <div v-else class="flex items-center gap-1.5 mt-1">
                <span class="w-2.5 h-2.5 rounded-full bg-slate-300 dark:bg-slate-600"></span>
                <span class="text-sm font-semibold text-slate-500 dark:text-slate-400">Inactivo</span>
              </div>
            </div>
          </div>
        </header>

        <!-- Recovered points notice -->
        <div v-if="recoveredPoints > 0" class="glass-card p-4 border-brand-200 dark:border-brand-500/30 bg-brand-50 dark:bg-brand-900/20 flex items-center gap-4 animate-in">
          <div class="p-2 bg-white dark:bg-brand-500/20 rounded-xl text-brand-500">
            <CloudArrowUpIcon class="w-6 h-6 flex-shrink-0" />
          </div>
          <div>
            <h4 class="text-brand-800 dark:text-brand-300 font-bold text-sm">Sincronización Completada</h4>
            <p class="text-brand-600/80 dark:text-brand-200/70 text-xs mt-0.5">Se recuperaron {{ recoveredPoints }} puntos GPS de tu sesión anterior.</p>
          </div>
        </div>

        <!-- Main Status Card -->
        <section class="glass-card overflow-hidden animate-slide-up" style="animation-delay: 0.1s;">
          <div class="p-6 sm:p-8 flex flex-col md:flex-row items-center gap-6 md:gap-10">
            
            <!-- Status Indicator -->
            <div class="relative group">
              <div class="absolute -inset-1 bg-gradient-to-r blur opacity-20 group-hover:opacity-40 transition duration-1000 rounded-full"
                :class="activeRecordId ? 'from-emerald-400 to-teal-500' : 'from-slate-400 to-slate-300 dark:from-slate-600 dark:to-slate-700'"></div>
              <div class="relative w-24 h-24 sm:w-32 sm:h-32 rounded-full flex items-center justify-center ring-4 ring-white dark:ring-surface-900 shadow-md transition-all"
                :class="activeRecordId ? 'bg-emerald-50 text-emerald-500 dark:bg-emerald-500/10 dark:text-emerald-400' : 'bg-slate-50 text-slate-400 dark:bg-surface-800 dark:text-slate-500'">
                <SignalIcon v-if="activeRecordId" class="w-10 h-10 sm:w-14 sm:h-14 animate-pulse-slow" />
                <MinusCircleIcon v-else class="w-10 h-10 sm:w-14 sm:h-14" />
              </div>
            </div>

            <!-- Status Details -->
            <div class="flex-1 text-center md:text-left">
              <h3 class="text-lg sm:text-xl font-bold tracking-tight mb-2"
                :class="activeRecordId ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-500 dark:text-slate-400'">
                {{ activeRecordId ? 'Jornada Activa' : 'No has iniciado turno' }}
              </h3>
              
              <div class="bg-slate-50 dark:bg-white/5 rounded-2xl p-4 inline-flex flex-col md:block w-full text-left">
                <div class="flex items-center gap-3">
                  <ClockIcon class="w-5 h-5 text-slate-400" />
                  <div>
                    <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Hora de Entrada</p>
                    <p class="text-sm font-bold text-slate-900 dark:text-white">{{ entryTime || '--:--' }}</p>
                  </div>
                </div>
                <!-- Additional Context if active -->
                <div v-if="activeRecordId" class="mt-3 pt-3 border-t border-slate-200 dark:border-white/10 flex items-center gap-3">
                  <MapPinIcon class="w-5 h-5 text-emerald-500" />
                  <div>
                    <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Puntos registrados</p>
                    <p class="text-sm font-bold text-slate-900 dark:text-white">{{ locationPoints.length }} pts</p>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- CTAs inside card on Desktop -->
            <div class="w-full md:w-auto flex-shrink-0">
              <button
                v-if="!activeRecordId"
                @click="startCheckProcess('entry')"
                :disabled="processing"
                class="btn-success btn-xl w-full md:min-w-[200px]"
              >
                <span v-if="processing" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <MapPinIcon v-else class="w-6 h-6" />
                <span>Registrar Entrada</span>
              </button>

              <button
                v-if="activeRecordId"
                @click="startCheckProcess('exit')"
                :disabled="processing"
                class="btn-danger btn-xl w-full md:min-w-[200px]"
              >
                <span v-if="processing" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <ArrowRightOnRectangleIcon v-else class="w-6 h-6" />
                <span>Registrar Salida</span>
              </button>
            </div>
          </div>
        </section>

      </div>

      <!-- Mobile Bottom Navigation Bar (Hidden on md up, Only shown on PWA standalone) -->
      <nav v-if="isPWA" class="md:hidden glass-panel fixed bottom-4 left-4 right-4 z-20 flex justify-between items-center px-6 py-3 border border-white/20 shadow-glass">
        <router-link to="/checkin" class="flex flex-col items-center gap-1 text-brand-600 dark:text-brand-400">
          <MapPinIcon class="w-6 h-6" />
          <span class="text-[10px] font-bold">Inicio</span>
        </router-link>
        
        <router-link to="/history" class="flex flex-col items-center gap-1 text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors">
          <ClockIcon class="w-6 h-6" />
          <span class="text-[10px] font-semibold">Historial</span>
        </router-link>

        <div class="relative -top-6">
          <router-link to="/qr" class="w-14 h-14 bg-gradient-to-br from-brand-500 to-violet-600 rounded-full flex items-center justify-center text-white shadow-glow-brand ring-4 ring-white dark:ring-surface-950 transition-transform active:scale-95">
            <QrCodeIcon class="w-6 h-6" />
          </router-link>
        </div>

        <router-link to="/profile" class="flex flex-col items-center gap-1 text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors">
          <UserCircleIcon class="w-6 h-6" />
          <span class="text-[10px] font-semibold">Perfil</span>
        </router-link>
        
        <button @click="showLogoutModal = true" class="flex flex-col items-center gap-1 text-slate-400 hover:text-rose-500 transition-colors">
          <ArrowRightOnRectangleIcon class="w-6 h-6" />
          <span class="text-[10px] font-semibold">Salir</span>
        </button>
      </nav>
      <!-- End mobile nav padding -->
      <div v-if="isPWA" class="h-24 md:h-0"></div>
    </main>

    <!-- ===== Modals ===== -->
    <Teleport to="body">

      <!-- Mobile Hamburger Menu Drawer -->
      <Transition name="slide-banner">
        <div v-if="showMobileMenu" class="fixed inset-0 z-50 flex justify-end">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showMobileMenu = false"></div>
          <div class="relative w-64 max-w-[80vw] h-full bg-white dark:bg-surface-950 shadow-2xl flex flex-col animate-slide-up sm:animate-in border-l border-slate-200 dark:border-surface-800">
            <div class="p-4 flex justify-between items-center border-b border-slate-100 dark:border-surface-800">
              <span class="font-bold text-slate-900 dark:text-white">Menú</span>
              <button @click="showMobileMenu = false" class="p-2 bg-slate-100 dark:bg-surface-800 rounded-full text-slate-500 dark:text-slate-400">
                <XMarkIcon class="w-5 h-5"/>
              </button>
            </div>
            <nav class="p-4 space-y-2 flex-1 overflow-y-auto">
               <router-link to="/checkin" @click="showMobileMenu = false" class="flex items-center gap-3 px-4 py-3 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400 font-semibold transition-all">
                  <MapPinIcon class="w-5 h-5" />
                  Inicio
                </router-link>
                <router-link to="/history" @click="showMobileMenu = false" class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
                  <ClockIcon class="w-5 h-5" />
                  Historial
                </router-link>
                <router-link to="/qr" @click="showMobileMenu = false" class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
                  <QrCodeIcon class="w-5 h-5" />
                  Escanear QR
                </router-link>
            </nav>
            <div class="p-4 border-t border-slate-100 dark:border-surface-800 flex gap-2">
               <button @click="showLogoutModal = true; showMobileMenu = false" class="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-rose-500 dark:text-rose-400 font-semibold hover:bg-rose-50 dark:hover:bg-rose-500/10 transition-colors">
                 <ArrowRightOnRectangleIcon class="w-5 h-5"/>
                 Cerrar Sesión
               </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Location error / guide modal -->
      <Transition name="modal">
        <div v-if="showLocationErrorModal" class="fixed inset-0 z-[60] flex items-end sm:items-center justify-center sm:p-4 bg-slate-900/60 dark:bg-slate-950/80 backdrop-blur-sm">
          <div class="w-full max-w-md glass-panel p-6 sm:p-8 animate-slide-up rounded-t-[2rem] sm:rounded-3xl border-b-0 sm:border-b border-t border-x border-white/20 dark:border-white/10 m-0">
            <div class="flex items-center justify-between mb-6">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 rounded-2xl flex items-center justify-center flex-shrink-0 bg-rose-50 dark:bg-rose-500/10 border border-rose-100 dark:border-rose-500/20 text-rose-500">
                  <MapPinIcon class="w-6 h-6" />
                </div>
                <div>
                  <h3 class="font-bold text-lg text-slate-900 dark:text-white">
                    {{ locationErrorType === 'inapp' ? 'Abre la página en Safari' :
                       locationErrorType === 'unavailable' ? 'GPS no disponible' :
                       locationErrorType === 'timeout' ? 'Tiempo de espera agotado' :
                       'Ubicación bloqueada' }}
                  </h3>
                  <p class="text-sm text-slate-500 dark:text-slate-400 mt-0.5">
                    {{ locationErrorType === 'inapp' ? 'El browser interno bloquea el GPS' :
                       locationErrorType === 'unavailable' ? 'No se detectó señal GPS' :
                       locationErrorType === 'timeout' ? 'El GPS tardó demasiado' :
                       'Sigue los pasos para activarla' }}
                    <span v-if="locationErrorCode" class="opacity-50">(Error: {{ locationErrorCode }})</span>
                  </p>
                </div>
              </div>
            </div>

            <!-- In-app browser (WKWebView) -->
            <template v-if="locationErrorType === 'inapp'">
              <div class="bg-brand-50 dark:bg-brand-500/5 border border-brand-100 dark:border-brand-500/20 rounded-2xl p-5 mb-6">
                <p class="text-sm font-bold mb-3 text-brand-600 dark:text-brand-400 uppercase tracking-wider">Cómo abrir en Safari</p>
                <ol class="space-y-3">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span class="w-6 h-6 rounded-full bg-brand-100 dark:bg-brand-500/20 text-brand-600 dark:text-brand-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">1</span>
                    <span class="mt-0.5">Toca el ícono de <strong class="text-slate-900 dark:text-white">tres puntos (···)</strong> en la barra inferior.</span>
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span class="w-6 h-6 rounded-full bg-brand-100 dark:bg-brand-500/20 text-brand-600 dark:text-brand-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">2</span>
                    <span class="mt-0.5">Selecciona <strong class="text-slate-900 dark:text-white">"Abrir en navegador"</strong>.</span>
                  </li>
                </ol>
              </div>
            </template>

            <!-- GPS unavailable or timeout -->
            <template v-else-if="locationErrorType === 'unavailable' || locationErrorType === 'timeout'">
              <div class="bg-amber-50 dark:bg-amber-500/5 border border-amber-100 dark:border-amber-500/20 rounded-2xl p-5 mb-6">
                <ol class="space-y-3">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span class="w-6 h-6 rounded-full bg-amber-100 dark:bg-amber-500/20 text-amber-600 dark:text-amber-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">1</span>
                    <span class="mt-0.5">Ve a un lugar con buena señal o al aire libre.</span>
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span class="w-6 h-6 rounded-full bg-amber-100 dark:bg-amber-500/20 text-amber-600 dark:text-amber-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">2</span>
                    <span class="mt-0.5">Asegúrate de que el GPS esté <strong class="text-slate-900 dark:text-white">activado</strong>.</span>
                  </li>
                </ol>
              </div>
            </template>

            <!-- Permission denied -->
            <template v-else>
              <!-- Paso 1: Configuración iOS global -->
              <div class="bg-brand-50 dark:bg-brand-500/5 border border-brand-100 dark:border-brand-500/20 rounded-2xl p-4 mb-4">
                <p class="text-xs font-bold mb-2.5 text-brand-600 dark:text-brand-400 uppercase tracking-wider">① Config. del iPhone (Revisar primero)</p>
                <ol class="space-y-2">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-brand-500 font-bold mt-px">1.</span> Ve a <strong>Configuración → Privacidad → Localización</strong>.
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-brand-500 font-bold mt-px">2.</span> Busca <strong>Safari</strong> y ponlo en <strong>"Al usar la app"</strong>.
                  </li>
                </ol>
              </div>

              <!-- Paso 2: AA del sitio en Safari -->
              <div class="bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-2xl p-4 mb-6">
                <p class="text-xs font-bold mb-2.5 text-slate-500 dark:text-slate-400 uppercase tracking-wider">② En esta página (Safari)</p>
                <ol class="space-y-2">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-slate-400 dark:text-slate-500 font-bold mt-px">1.</span> Toca las letras <strong>"AA"</strong> en la barra web.
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-slate-400 dark:text-slate-500 font-bold mt-px">2.</span> Ajustes del sitio web → Ubicación <strong>"Permitir"</strong>.
                  </li>
                </ol>
              </div>
            </template>

            <div class="space-y-3 sm:space-y-0 sm:grid sm:grid-cols-2 sm:gap-3">
              <button @click="cancelProcess" class="btn-secondary btn-lg w-full order-last sm:order-first">Cancerlar</button>
              <button v-if="locationErrorType === 'denied'" @click="reloadPage()" class="btn-primary btn-lg w-full">
                <ArrowPathRoundedSquareIcon class="w-5 h-5" /> Recargar
              </button>
              <button v-else-if="locationErrorType !== 'inapp'" @click="retryLocation()" class="btn-primary btn-lg w-full">
                <ArrowPathRoundedSquareIcon class="w-5 h-5" /> Reintentar
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Camera modal -->
      <div v-if="showCameraModal" class="fixed inset-0 bg-black z-50 flex flex-col">
        <!-- Camera instruction bar -->
        <div class="relative z-10 text-center py-5 px-6 safe-top"
          style="background: linear-gradient(to bottom, rgba(0,0,0,0.85), transparent);">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full text-white text-sm font-semibold"
            style="background: rgba(99,102,241,0.5); border: 1px solid rgba(99,102,241,0.5);">
            <CameraIcon v-if="cameraStep === 'site'" class="w-4 h-4" />
            <UserCircleIcon v-else class="w-4 h-4" />
            <span>{{ cameraStep === 'site' ? 'Toma fotografía del sitio' : 'Tómate una selfie en el sitio' }}</span>
          </div>
          <p class="text-white/60 text-xs mt-2">
            {{ cameraStep === 'site' ? 'Fotografía el lugar donde te encuentras' : 'Incluye claramente tu rostro' }}
          </p>
        </div>

        <!-- Video / Preview -->
        <div class="flex-1 relative overflow-hidden">
          <video ref="videoRef" autoplay playsinline muted
            class="absolute inset-0 w-full h-full object-cover"
            :class="{ 'scale-x-[-1]': facingFront }"></video>
          <canvas ref="canvasRef" class="hidden"></canvas>

          <!-- Photo preview overlay -->
          <div v-if="capturedPhoto" class="absolute inset-0 bg-black">
            <img :src="capturedPhoto" class="w-full h-full object-cover"
              :class="{ 'scale-x-[-1]': cameraStep === 'selfie' }" />
            <!-- Preview overlay badge -->
            <div class="absolute top-4 left-1/2 -translate-x-1/2 px-4 py-2 rounded-full text-sm font-semibold"
              style="background: rgba(0,0,0,0.7); border: 1px solid rgba(255,255,255,0.2); color: white;">
              Vista previa
            </div>
          </div>

          <!-- Guide frame -->
          <div v-if="!capturedPhoto && cameraStep === 'selfie'"
            class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div class="w-48 h-60 rounded-3xl border-2 border-white/30 border-dashed"></div>
          </div>
        </div>

        <!-- Controls bar -->
        <div class="safe-bottom py-6 px-6"
          style="background: linear-gradient(to top, rgba(0,0,0,0.9), transparent 80%);">

          <!-- Before capture -->
          <div v-if="!capturedPhoto" class="flex items-center justify-between max-w-xs mx-auto">
            <button @click="toggleFlash"
              class="w-12 h-12 rounded-full flex items-center justify-center text-2xl transition-all"
              :style="flashOn ? 'background: rgba(251,191,36,0.3); border: 2px solid rgba(251,191,36,0.7);' : 'background: rgba(255,255,255,0.1); border: 2px solid rgba(255,255,255,0.2);'">
              <BoltIcon class="w-6 h-6" :class="flashOn ? 'text-amber-300' : 'text-white'" />
            </button>

            <button @click="capturePhoto"
              class="w-20 h-20 rounded-full flex items-center justify-center transition-all active:scale-90"
              style="background: white; border: 4px solid rgba(255,255,255,0.3); box-shadow: 0 0 0 6px rgba(255,255,255,0.1);">
              <div class="w-14 h-14 rounded-full border-2 border-gray-300" style="background: white;"></div>
            </button>

            <button @click="switchCamera"
              class="w-12 h-12 rounded-full flex items-center justify-center text-2xl transition-all"
              style="background: rgba(255,255,255,0.1); border: 2px solid rgba(255,255,255,0.2);">
              <ArrowPathIcon class="w-6 h-6 text-white" />
            </button>
          </div>

          <!-- After capture -->
          <div v-if="capturedPhoto" class="flex gap-3 max-w-xs mx-auto">
            <button @click="retakePhoto" class="btn-secondary btn-lg flex-1">
              <ArrowPathIcon class="w-4 h-4" /> Retomar
            </button>
            <button @click="confirmPhoto" class="btn-primary btn-lg flex-1">
              <CheckCircleIcon class="w-4 h-4" /> Usar foto
            </button>
          </div>
        </div>
      </div>

      <!-- Success modal -->
      <Transition name="modal">
        <div v-if="showSuccessModal" class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0,0,0,0.75); backdrop-filter: blur(8px);">
          <div class="w-full max-w-sm glass-card p-8 text-center animate-in">
            <div class="w-20 h-20 rounded-3xl flex items-center justify-center mx-auto mb-4"
              :style="checkType === 'entry'
                ? 'background: linear-gradient(135deg, rgba(5,150,105,0.3), rgba(16,185,129,0.3)); border: 1px solid rgba(16,185,129,0.4);'
                : 'background: linear-gradient(135deg, rgba(220,38,38,0.3), rgba(244,63,94,0.3)); border: 1px solid rgba(244,63,94,0.4);'">
              <CheckCircleIcon v-if="checkType === 'entry'" class="w-10 h-10 text-emerald-400" />
              <HandRaisedIcon v-else class="w-10 h-10 text-rose-400" />
            </div>
            <h3 class="text-2xl font-bold mb-1" style="color: var(--text);">
              {{ checkType === 'entry' ? '¡Entrada Registrada!' : '¡Salida Registrada!' }}
            </h3>
            <p style="color: var(--text-muted);" class="text-sm mb-1">{{ registeredAt }}</p>
            <p class="text-xs mb-6" style="color: var(--text-muted);">Ubicación y fotografías guardadas exitosamente</p>
            <button @click="closeSuccess" class="btn-primary btn-lg w-full">Aceptar</button>
          </div>
        </div>
      </Transition>

      <!-- Processing spinner -->
      <div v-if="processing && !showLocationErrorModal && !showCameraModal && !showSuccessModal"
        class="fixed inset-0 z-50 flex items-center justify-center"
        style="background: rgba(0,0,0,0.5); backdrop-filter: blur(4px);">
        <div class="glass-card p-6 text-center">
          <div class="w-10 h-10 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin mx-auto mb-3"></div>
          <p class="font-semibold" style="color: var(--text);">Procesando...</p>
          <p class="text-xs mt-1" style="color: var(--text-muted);">Por favor espera</p>
        </div>
      </div>

      <!-- Logout confirmation modal -->
      <Transition name="modal">
        <div v-if="showLogoutModal" class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0,0,0,0.85); backdrop-filter: blur(8px);">
          <div class="w-full max-w-sm glass-card p-6 animate-in text-center" style="background: var(--modal-bg);">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center bg-rose-500/10 text-rose-500">
              <ArrowRightOnRectangleIcon class="w-8 h-8" />
            </div>
            <h3 class="font-bold text-lg mb-2" style="color: var(--text);">¿Cerrar sesión?</h3>
            <p class="text-sm mb-6" style="color: var(--text-muted);">Tendrás que volver a ingresar tus credenciales para registrar asistencia.</p>
            <div class="flex gap-3">
              <button @click="showLogoutModal = false" class="btn-secondary flex-1">Cancelar</button>
              <button @click="confirmLogout" class="btn-danger flex-1">Salir</button>
            </div>
          </div>
        </div>
      </Transition>

    </Teleport>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import {
  CheckCircleIcon, HandRaisedIcon, MapPinIcon, ArrowRightOnRectangleIcon,
  ExclamationTriangleIcon, XMarkIcon, LockClosedIcon, CloudArrowUpIcon,
  SignalIcon, MinusCircleIcon, CameraIcon, UserCircleIcon, BoltIcon, ArrowPathIcon,
  InformationCircleIcon, DevicePhoneMobileIcon, ArrowPathRoundedSquareIcon,
  QrCodeIcon, ClockIcon, Bars3BottomRightIcon
} from '@heroicons/vue/24/outline'

const auth = useAuthStore()
const router = useRouter()

const showLogoutModal = ref(false)

async function confirmLogout() {
  showLogoutModal.value = false
  await auth.logout()
  router.push('/login')
}

// Clock
const currentTime = ref('')
const todayStr = ref('')
let timeInterval = null

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  todayStr.value = now.toLocaleDateString('es-MX', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}


// State
const activeRecordId = ref(localStorage.getItem('activeRecordId') || null)
const entryTime = ref(localStorage.getItem('entryTime') || null)
const processing = ref(false)
const checkType = ref('entry')
const showLocationErrorModal = ref(false)
const locationErrorType = ref('denied')
const locationErrorCode = ref(0)
const showCameraModal = ref(false)
const showSuccessModal = ref(false)
const registeredAt = ref('')

const isPWA = ref(false)
const showMobileMenu = ref(false)

onMounted(() => {
  // Check if PWA (standalone mode)
  if (window.matchMedia('(display-mode: standalone)').matches || window.navigator.standalone) {
    isPWA.value = true
  }
  window.matchMedia('(display-mode: standalone)').addEventListener('change', (e) => {
    isPWA.value = e.matches
  })
})

// Location
let locationWatchId = null
const locationPoints = []

// UI extras
const recoveredPoints = ref(0)
const showHiddenWarning = ref(false)
const isTracking = ref(false)

// Wake lock
let wakeLock = null
const wakeLockActive = ref(false)
async function requestWakeLock() {
  if ('wakeLock' in navigator) {
    try {
      wakeLock = await navigator.wakeLock.request('screen')
      wakeLockActive.value = true
      wakeLock.addEventListener('release', () => { wakeLockActive.value = false })
    } catch { /* not supported */ }
  }
}
function releaseWakeLock() {
  if (wakeLock) { wakeLock.release(); wakeLock = null; wakeLockActive.value = false }
}

// Location buffer — persists points when app goes to background
function saveLocationBuffer(points) {
  try { localStorage.setItem('locationBuffer', JSON.stringify(points)) } catch {}
}
function loadLocationBuffer() {
  try { return JSON.parse(localStorage.getItem('locationBuffer') || '[]') } catch { return [] }
}
function clearLocationBuffer() {
  try { localStorage.removeItem('locationBuffer') } catch {}
}
function handleVisibilityChange() {
  if (document.visibilityState === 'hidden') {
    saveLocationBuffer(locationPoints)
  } else if (document.visibilityState === 'visible' && isTracking.value) {
    showHiddenWarning.value = true
  }
}

function reloadPage() { window.location.reload() }

function isIOS() {
  return /iPhone|iPad|iPod/.test(navigator.userAgent)
}

function isIOSInAppBrowser() {
  const ua = navigator.userAgent
  if (!isIOS()) return false
  // Pure WKWebView (older WhatsApp, Telegram, etc.) — no Version/ token
  if (!/Version\//.test(ua)) return true
  // SFSafariViewController shares the Safari UA but lacks the actual Safari window object.
  // Detect by checking standalone is unavailable AND no window.safari (removed iOS 16 but still works on 15-)
  // Best heuristic: check known in-app UA tokens
  if (/FBAN|FBAV|Instagram|Twitter|LinkedInApp|Snapchat/.test(ua)) return true
  return false
}

function startCheckProcess(type) {
  checkType.value = type
  if (isIOSInAppBrowser()) {
    locationErrorType.value = 'inapp'
    showLocationErrorModal.value = true
    return
  }
  // Call requestLocation directly — iOS Safari requires getCurrentPosition
  // to be called in the same synchronous call stack as the user tap.
  requestLocation()
}

function cancelProcess() {
  showLocationErrorModal.value = false
  processing.value = false
}

// requestLocation MUST be a plain synchronous function so that
// navigator.geolocation.getCurrentPosition() is called as the VERY FIRST
// instruction — before any Vue reactive state changes — to stay within the
// iOS Safari user-gesture call stack. Even a single ref assignment before
// the call can cause WebKit to break the gesture chain.
function requestLocation() {
  if (!navigator.geolocation) {
    locationErrorCode.value = 0
    locationErrorType.value = 'unavailable'
    showLocationErrorModal.value = true
    return
  }
  // ← getCurrentPosition is the absolute first call. No state changes before this.
  navigator.geolocation.getCurrentPosition(
    async (position) => {
      processing.value = false
      locationPoints.push({
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
        accuracy: position.coords.accuracy,
        recorded_at: new Date().toISOString()
      })
      locationWatchId = navigator.geolocation.watchPosition(
        pos => locationPoints.push({
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
          accuracy: pos.coords.accuracy,
          recorded_at: new Date().toISOString()
        }),
        null,
        { enableHighAccuracy: true, maximumAge: 5000 }
      )
      isTracking.value = true
      requestWakeLock()
      cameraStep.value = 'site'
      showCameraModal.value = true
      await nextTick()
      await startCamera('environment')
    },
    (err) => {
      processing.value = false
      const code = err?.code ?? -1
      locationErrorCode.value = code
      if (code === 1) locationErrorType.value = 'denied'           // PERMISSION_DENIED
      else if (code === 2) locationErrorType.value = 'unavailable'  // POSITION_UNAVAILABLE
      else locationErrorType.value = 'timeout'                      // TIMEOUT or unknown
      showLocationErrorModal.value = true
    },
    { enableHighAccuracy: true, timeout: 30000, maximumAge: 0 }
  )
  // Show processing spinner after registering the geolocation listener
  processing.value = true
}

function stopLocationTracking() {
  if (locationWatchId !== null) {
    navigator.geolocation.clearWatch(locationWatchId)
    locationWatchId = null
  }
  isTracking.value = false
  releaseWakeLock()
}

function retryLocation() {
  showLocationErrorModal.value = false
  requestLocation()
}

// Camera
const videoRef = ref(null)
const canvasRef = ref(null)
const capturedPhoto = ref(null)
const cameraStep = ref('site')
const facingFront = ref(false)
const flashOn = ref(false)
let stream = null
let photoSite = null
let photoSelfie = null

async function startCamera(facing = 'environment') {
  stopCamera()
  facingFront.value = facing === 'user'
  stream = await navigator.mediaDevices.getUserMedia({
    video: { facingMode: facing, width: { ideal: 1280 }, height: { ideal: 720 } },
    audio: false
  })
  if (videoRef.value) videoRef.value.srcObject = stream
}

function stopCamera() {
  if (stream) { stream.getTracks().forEach(t => t.stop()); stream = null }
}

async function switchCamera() {
  await startCamera(facingFront.value ? 'environment' : 'user')
}

function toggleFlash() {
  const track = stream?.getVideoTracks()[0]
  if (!track) return
  const caps = track.getCapabilities?.()
  if (caps?.torch) {
    flashOn.value = !flashOn.value
    track.applyConstraints({ advanced: [{ torch: flashOn.value }] })
  } else {
    alert('Tu dispositivo no soporta flash desde el navegador')
  }
}

function capturePhoto() {
  const video = videoRef.value
  const canvas = canvasRef.value
  if (!video || !canvas) return

  // Limit resolution to max 800px to keep payload under 1 MB
  const MAX_DIM = 800
  const scale = Math.min(1, MAX_DIM / Math.max(video.videoWidth, video.videoHeight))
  canvas.width = Math.round(video.videoWidth * scale)
  canvas.height = Math.round(video.videoHeight * scale)

  const ctx = canvas.getContext('2d')
  if (cameraStep.value === 'selfie') { ctx.translate(canvas.width, 0); ctx.scale(-1, 1) }
  ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
  capturedPhoto.value = canvas.toDataURL('image/jpeg', 0.72)
}

function retakePhoto() { capturedPhoto.value = null }

async function confirmPhoto() {
  if (cameraStep.value === 'site') {
    photoSite = capturedPhoto.value
    capturedPhoto.value = null
    cameraStep.value = 'selfie'
    await startCamera('user')
  } else {
    photoSelfie = capturedPhoto.value
    capturedPhoto.value = null
    showCameraModal.value = false
    stopCamera()
    await submitCheck()
  }
}

function dataURLtoBlob(dataURL) {
  const [header, data] = dataURL.split(',')
  const mime = header.match(/:(.*?);/)[1]
  const bstr = atob(data)
  const u8arr = new Uint8Array(bstr.length)
  for (let i = 0; i < bstr.length; i++) u8arr[i] = bstr.charCodeAt(i)
  return new Blob([u8arr], { type: mime })
}

async function submitCheck() {
  processing.value = true
  try {
    const now = new Date().toISOString()
    const formData = new FormData()
    formData.append('type', checkType.value)
    formData.append('timestamp', now)
    // Send GPS coords so backend can cross-check against IP geolocation
    if (locationPoints.length > 0) {
      formData.append('latitude', String(locationPoints[0].latitude))
      formData.append('longitude', String(locationPoints[0].longitude))
    }
    if (photoSite) formData.append('photo_site', dataURLtoBlob(photoSite), 'site.jpg')
    if (photoSelfie) formData.append('photo_selfie', dataURLtoBlob(photoSelfie), 'selfie.jpg')

    const { data } = await api.post('/checks', formData)
    const recordId = data.record.id

    if (locationPoints.length > 0) {
      try {
        await api.post('/location-points/batch', {
          check_record_id: recordId,
          points: locationPoints
        })
      } catch { /* Points failed — non-critical */ }
    }
    locationPoints.length = 0
    clearLocationBuffer()
    stopLocationTracking()
    recoveredPoints.value = 0

    if (checkType.value === 'entry') {
      activeRecordId.value = recordId
      entryTime.value = new Date(now).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
      localStorage.setItem('activeRecordId', recordId)
      localStorage.setItem('entryTime', entryTime.value)
    } else {
      activeRecordId.value = null
      entryTime.value = null
      localStorage.removeItem('activeRecordId')
      localStorage.removeItem('entryTime')
    }

    registeredAt.value = new Date(now).toLocaleString('es-MX')
    showSuccessModal.value = true
  } catch (err) {
    const status = err?.response?.status
    const msg = err?.response?.data?.error || ''
    alert(`Error al registrar. Intenta de nuevo.${status ? ` (${status}${msg ? ': ' + msg : ''})` : ''}`)
  } finally {
    processing.value = false
    photoSite = null
    photoSelfie = null
  }
}

function closeSuccess() { showSuccessModal.value = false }

// ─── Lifecycle ────────────────────────────────────────────────────────────────
onMounted(async () => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  document.addEventListener('visibilitychange', handleVisibilityChange)

  // Sync state with backend so it survives cache clearing
  try {
    const { data } = await api.get('/checks/status')
    if (data && data.active) {
      activeRecordId.value = data.record_id
      entryTime.value = new Date(data.entry_time).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
      localStorage.setItem('activeRecordId', data.record_id)
      localStorage.setItem('entryTime', entryTime.value)
    } else {
      activeRecordId.value = null
      entryTime.value = null
      localStorage.removeItem('activeRecordId')
      localStorage.removeItem('entryTime')
    }
  } catch (error) {
    console.error('Failed to sync session status', error)
  }

  // Recover buffered points if app was closed during an active session
  const buffered = loadLocationBuffer()
  if (buffered.length > 0 && activeRecordId.value) {
    recoveredPoints.value = buffered.length
    locationPoints.push(...buffered)
  }
})

onUnmounted(() => {
  clearInterval(timeInterval)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  stopLocationTracking()
  stopCamera()
})
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.25s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.safe-top { padding-top: max(16px, env(safe-area-inset-top)); }
.safe-bottom { padding-bottom: max(24px, env(safe-area-inset-bottom)); }
.slide-banner-enter-active, .slide-banner-leave-active { transition: all 0.3s ease; }
.slide-banner-enter-from, .slide-banner-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
