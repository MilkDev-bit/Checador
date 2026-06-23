<template>
  <div class="min-h-screen min-h-dvh bg-mesh flex flex-col md:flex-row">

    <!-- Sidebar / Top Navigation (Responsive) -->
    <aside
      class="md:w-64 glass-panel border-b md:border-b-0 md:border-r border-white/20 dark:border-white/5 flex-shrink-0 flex flex-col z-20 sticky top-0 md:h-screen md:sticky rounded-none">

      <!-- Brand & Header -->
      <div
        class="px-4 md:px-6 py-3 min-h-[4rem] safe-top flex items-center justify-between border-b border-slate-200 dark:border-white/5">
        <div class="flex items-center gap-3">
          <div
            class="w-9 h-9 rounded-xl overflow-hidden flex items-center justify-center flex-shrink-0 bg-gradient-to-br from-brand-500 to-violet-600 shadow-glow-brand ring-2 ring-white/20 dark:ring-white/5">
            <img src="/LOGOC.png" alt="Logo" class="w-full h-full object-contain" />
          </div>
          <div class="flex flex-col justify-center">
            <h1 class="font-bold text-base leading-tight tracking-tight text-slate-900 dark:text-white">PaseLista</h1>
            <p class="text-[11px] font-medium text-slate-500 dark:text-slate-400 truncate max-w-[120px]">{{
              auth.user?.project_name || 'Mi Proyecto' }}</p>
          </div>
        </div>

        <!-- Mobile Actions -->
        <div class="flex items-center gap-2 md:hidden">
          <router-link to="/profile"
            class="w-8 h-8 rounded-full overflow-hidden flex-shrink-0 transition-transform active:scale-95 ring-2 ring-brand-100 dark:ring-brand-900/50">
            <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" class="w-full h-full object-cover"
              alt="Perfil" />
            <div v-else
              class="w-full h-full bg-brand-500 flex items-center justify-center text-white text-xs font-bold">
              {{ auth.user?.first_name?.[0] ?? '' }}{{ auth.user?.last_name?.[0] ?? '' }}
            </div>
          </router-link>
          <button v-if="!isPWA" @click="showMobileMenu = true"
            class="w-8 h-8 ml-2 flex items-center justify-center text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-surface-800 rounded-lg transition-colors">
            <Bars3BottomRightIcon class="w-6 h-6" />
          </button>
        </div>
      </div>

      <!-- Navigation Links (Hidden on small screens, shown on md up) -->
      <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-1.5 hidden md:block">
        <p class="px-3 text-[11px] font-bold uppercase tracking-wider text-slate-400 dark:text-slate-500 mb-3">Menú
          Principal</p>

        <router-link to="/checkin"
          class="flex items-center gap-3 px-3 py-2.5 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400 font-semibold group transition-all">
          <MapPinIcon class="w-5 h-5 group-hover:scale-110 transition-transform" />
          Mi Check-in
        </router-link>

        <router-link to="/history"
          class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 hover:text-slate-900 dark:hover:text-white group transition-all">
          <ClockIcon class="w-5 h-5 group-hover:scale-110 transition-transform" />
          Historial
        </router-link>

        <router-link to="/qr"
          class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 hover:text-slate-900 dark:hover:text-white group transition-all">
          <QrCodeIcon class="w-5 h-5 group-hover:scale-110 transition-transform" />
          Escanear QR
        </router-link>
      </nav>

      <!-- Desktop Bottom Actions -->
      <div class="p-4 border-t border-slate-200 dark:border-white/5 hidden md:block space-y-3">
        <router-link to="/profile"
          class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
          <div
            class="w-8 h-8 rounded-full overflow-hidden flex-shrink-0 bg-brand-100 dark:bg-brand-900/40 border border-brand-200 dark:border-brand-800">
            <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" class="w-full h-full object-cover"
              alt="Perfil" />
            <div v-else
              class="w-full h-full text-brand-600 dark:text-brand-400 flex items-center justify-center text-xs font-bold">
              {{ auth.user?.first_name?.[0] ?? '' }}
            </div>
          </div>
          <div class="overflow-hidden">
            <p class="text-sm font-semibold text-slate-900 dark:text-white truncate">{{ auth.user?.first_name }} {{
              auth.user?.last_name }}</p>
            <p class="text-xs text-slate-500 dark:text-slate-400 truncate">Ver perfil</p>
          </div>
        </router-link>

        <button @click="showLogoutModal = true"
          class="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-rose-500 dark:text-rose-400 font-semibold hover:bg-rose-50 dark:hover:bg-rose-500/10 transition-colors mt-2">
          <ArrowLeftOnRectangleIcon class="w-4 h-4" />
          Cerrar Sesión
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 relative overflow-y-auto flex flex-col md:pt-6 md:px-8 max-w-5xl mx-auto w-full">

      <!-- GPS background warning banner -->
      <Transition name="slide-banner">
        <div v-if="showHiddenWarning"
          class="glass-card mb-4 mx-4 md:mx-0 p-4 border border-amber-200 dark:border-amber-500/20 bg-amber-50/80 dark:bg-amber-900/20 flex items-start sm:items-center gap-4">
          <div class="p-2 bg-amber-100 dark:bg-amber-500/20 rounded-full text-amber-500 dark:text-amber-400">
            <ExclamationTriangleIcon class="w-5 h-5 flex-shrink-0" />
          </div>
          <div class="flex-1">
            <h4 class="text-amber-800 dark:text-amber-400 font-bold text-sm">GPS posiblemente pausado</h4>
            <p class="text-amber-700/80 dark:text-amber-300/70 text-xs mt-0.5 leading-relaxed">La app estuvo en segundo
              plano. Algunos puntos podrían no haberse registrado.</p>
          </div>
          <button @click="showHiddenWarning = false"
            class="p-1.5 hover:bg-amber-200/50 dark:hover:bg-amber-500/20 rounded-lg text-amber-600 dark:text-amber-400/60 transition-colors">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </Transition>

      <div class="px-4 py-6 md:py-0 space-y-6 flex-1">

        <!-- Header Section -->
        <header class="flex flex-col sm:flex-row sm:items-end justify-between gap-4 animate-slide-up">
          <div>
            <h2 class="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">Hola, {{
              auth.user?.first_name }} 👋</h2>
            <p class="text-slate-500 dark:text-slate-400 font-medium mt-1">{{ todayStr }}</p>
          </div>
          <div
            class="glass-card px-5 py-3 border border-slate-200 dark:border-white/10 flex items-center justify-between gap-6 shadow-sm">
            <div class="flex flex-col">
              <span class="text-[10px] uppercase font-bold tracking-wider text-slate-400 dark:text-slate-500">Hora
                actual</span>
              <span class="text-xl font-bold font-mono text-brand-600 dark:text-brand-400">{{ currentTime }}</span>
            </div>
            <div class="h-10 w-px bg-slate-200 dark:bg-white/10"></div>
            <div class="flex flex-col items-end">
              <span class="text-[10px] uppercase font-bold tracking-wider text-slate-400 dark:text-slate-500">Status
                GPS</span>
              <!-- Adquiriendo fix preciso -->
              <div v-if="gpsAcquiring" class="flex items-center gap-1.5 mt-1">
                <span class="relative flex h-2.5 w-2.5">
                  <span
                    class="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-amber-500"></span>
                </span>
                <span class="text-sm font-semibold text-amber-600 dark:text-amber-400">
                  {{ gpsAccuracy != null ? `~${gpsAccuracy}m` : 'Buscando…' }}
                </span>
              </div>
              <!-- Rastreando activamente -->
              <div v-else-if="isTracking" class="flex items-center gap-1.5 mt-1">
                <span class="relative flex h-2.5 w-2.5">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75"
                    :class="gpsAccuracy != null && gpsAccuracy <= 30 ? 'bg-emerald-400' : gpsAccuracy != null && gpsAccuracy <= 100 ? 'bg-amber-400' : 'bg-rose-400'"></span>
                  <span class="relative inline-flex rounded-full h-2.5 w-2.5"
                    :class="gpsAccuracy != null && gpsAccuracy <= 30 ? 'bg-emerald-500' : gpsAccuracy != null && gpsAccuracy <= 100 ? 'bg-amber-500' : 'bg-rose-500'"></span>
                </span>
                <span class="text-sm font-semibold"
                  :class="gpsAccuracy != null && gpsAccuracy <= 30 ? 'text-emerald-600 dark:text-emerald-400' : gpsAccuracy != null && gpsAccuracy <= 100 ? 'text-amber-600 dark:text-amber-400' : 'text-rose-500'">
                  {{ gpsAccuracy != null ? `±${gpsAccuracy}m` : 'Rastreando' }}
                </span>
              </div>
              <!-- Permiso concedido pero no rastreando -->
              <div v-else-if="gpsPermission === 'granted'" class="flex items-center gap-1.5 mt-1">
                <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
                <span class="text-sm font-semibold text-emerald-600 dark:text-emerald-400">Activo</span>
              </div>
              <!-- Permiso denegado -->
              <div v-else-if="gpsPermission === 'denied'" class="flex items-center gap-1.5 mt-1">
                <span class="w-2.5 h-2.5 rounded-full bg-rose-500"></span>
                <span class="text-sm font-semibold text-rose-500">Sin permiso</span>
              </div>
              <!-- Sin respuesta aún -->
              <div v-else class="flex items-center gap-1.5 mt-1">
                <span class="w-2.5 h-2.5 rounded-full bg-slate-300 dark:bg-slate-600"></span>
                <span class="text-sm font-semibold text-slate-500 dark:text-slate-400">Inactivo</span>
              </div>
            </div>
          </div>
        </header>

        <!-- Recovered points notice -->
        <div v-if="recoveredPoints > 0"
          class="glass-card p-4 border-brand-200 dark:border-brand-500/30 bg-brand-50 dark:bg-brand-900/20 flex items-center gap-4 animate-in">
          <div class="p-2 bg-white dark:bg-brand-500/20 rounded-xl text-brand-500">
            <CloudArrowUpIcon class="w-6 h-6 flex-shrink-0" />
          </div>
          <div>
            <h4 class="text-brand-800 dark:text-brand-300 font-bold text-sm">Sincronización Completada</h4>
            <p class="text-brand-600/80 dark:text-brand-200/70 text-xs mt-0.5">Se recuperaron {{ recoveredPoints }}
              puntos GPS de tu sesión anterior.</p>
          </div>
        </div>

        <!-- Pending offline sync notice -->
        <div v-if="pendingCount > 0"
          class="glass-card p-4 border-amber-200 dark:border-amber-500/30 bg-amber-50 dark:bg-amber-900/20 flex items-center gap-4 animate-in">
          <div class="p-2 bg-white dark:bg-amber-500/20 rounded-xl text-amber-500">
            <CloudArrowUpIcon class="w-6 h-6 flex-shrink-0 animate-bounce" />
          </div>
          <div class="flex-1">
            <h4 class="text-amber-800 dark:text-amber-300 font-bold text-sm">
              {{ pendingCount }} registro{{ pendingCount > 1 ? 's' : '' }} pendiente{{ pendingCount > 1 ? 's' : '' }} de
              sincronizar
            </h4>
            <p class="text-amber-600/80 dark:text-amber-200/70 text-xs mt-0.5">Se enviará{{ pendingCount > 1 ? 'n' : ''
              }} automáticamente al recuperar la conexión.</p>
          </div>
          <button v-if="isOnline" @click="syncPendingChecks"
            class="text-xs font-bold text-amber-700 dark:text-amber-300 px-3 py-1.5 rounded-lg bg-amber-100 dark:bg-amber-500/20 hover:bg-amber-200 dark:hover:bg-amber-500/30 transition-colors">
            Sincronizar
          </button>
        </div>

        <!-- Main Status Card -->
        <section class="glass-card overflow-hidden animate-slide-up" style="animation-delay: 0.1s;">

          <!-- Skeleton: loading state -->
          <div v-if="statusLoading"
            class="p-6 sm:p-8 flex flex-col md:flex-row items-center gap-6 md:gap-10 animate-pulse">
            <div class="w-24 h-24 sm:w-32 sm:h-32 rounded-full bg-slate-200 dark:bg-white/10 flex-shrink-0"></div>
            <div class="flex-1 w-full space-y-3">
              <div class="h-5 bg-slate-200 dark:bg-white/10 rounded-lg w-2/5 mx-auto md:mx-0"></div>
              <div class="h-16 bg-slate-200 dark:bg-white/10 rounded-2xl w-full"></div>
            </div>
            <div class="w-full md:w-48 h-12 bg-slate-200 dark:bg-white/10 rounded-xl flex-shrink-0"></div>
          </div>

          <!-- Error: no connection -->
          <div v-else-if="statusError" class="p-6 sm:p-8 flex flex-col items-center gap-4 text-center">
            <div
              class="w-16 h-16 rounded-2xl bg-rose-50 dark:bg-rose-500/10 flex items-center justify-center text-rose-500">
              <ExclamationTriangleIcon class="w-8 h-8" />
            </div>
            <div>
              <h3 class="font-bold text-slate-900 dark:text-white">Sin conexión al servidor</h3>
              <p class="text-sm text-slate-500 dark:text-slate-400 mt-1">Verifica tu internet para cargar el estado de
                tu sesión.</p>
            </div>
            <button @click="retryStatus" class="btn-primary btn-lg">
              <ArrowPathRoundedSquareIcon class="w-5 h-5" /> Reintentar
            </button>
          </div>

          <!-- Normal content -->
          <div v-else class="p-6 sm:p-8 flex flex-col md:flex-row items-center gap-6 md:gap-10">

            <!-- Status Indicator -->
            <div class="relative group">
              <div
                class="absolute -inset-1 bg-gradient-to-r blur opacity-20 group-hover:opacity-40 transition duration-1000 rounded-full"
                :class="activeRecordId ? 'from-emerald-400 to-teal-500' : 'from-slate-400 to-slate-300 dark:from-slate-600 dark:to-slate-700'">
              </div>
              <div
                class="relative w-24 h-24 sm:w-32 sm:h-32 rounded-full flex items-center justify-center ring-4 ring-white dark:ring-surface-900 shadow-md transition-all"
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
                <div v-if="activeRecordId"
                  class="mt-3 pt-3 border-t border-slate-200 dark:border-white/10 flex items-center gap-3">
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
              <button v-if="isSyncing" disabled
                class="btn-success btn-xl w-full md:min-w-[200px] opacity-70 cursor-not-allowed">
                <ArrowPathIcon class="w-6 h-6 animate-spin" />
                <span>Sincronizando...</span>
              </button>
              <button v-else-if="!activeRecordId" @click="startCheckProcess('entry')" :disabled="processing"
                class="btn-success btn-xl w-full md:min-w-[200px]">
                <span v-if="processing"
                  class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <MapPinIcon v-else class="w-6 h-6" />
                <span>Registrar Entrada</span>
              </button>

              <button v-if="activeRecordId" @click="startCheckProcess('exit')" :disabled="processing"
                class="btn-danger btn-xl w-full md:min-w-[200px]">
                <span v-if="processing"
                  class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <ArrowRightOnRectangleIcon v-else class="w-6 h-6" />
                <span>Registrar Salida</span>
              </button>
            </div>
          </div>
        </section>

      </div>

      <!-- Mobile Bottom Navigation Bar (Hidden on md up, Only shown on PWA standalone) -->
      <nav v-if="isPWA"
        class="md:hidden glass-panel fixed bottom-4 left-4 right-4 z-20 flex justify-between items-center px-6 py-3 border border-white/20 shadow-glass">
        <router-link to="/checkin" class="flex flex-col items-center gap-1 text-brand-600 dark:text-brand-400">
          <MapPinIcon class="w-6 h-6" />
          <span class="text-[10px] font-bold">Inicio</span>
        </router-link>

        <router-link to="/history"
          class="flex flex-col items-center gap-1 text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors">
          <ClockIcon class="w-6 h-6" />
          <span class="text-[10px] font-semibold">Historial</span>
        </router-link>

        <div class="relative -top-6">
          <router-link to="/qr"
            class="w-14 h-14 bg-gradient-to-br from-brand-500 to-violet-600 rounded-full flex items-center justify-center text-white shadow-glow-brand ring-4 ring-white dark:ring-surface-950 transition-transform active:scale-95">
            <QrCodeIcon class="w-6 h-6" />
          </router-link>
        </div>

        <router-link to="/profile"
          class="flex flex-col items-center gap-1 text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors">
          <UserCircleIcon class="w-6 h-6" />
          <span class="text-[10px] font-semibold">Perfil</span>
        </router-link>

        <button @click="showLogoutModal = true"
          class="flex flex-col items-center gap-1 text-slate-400 hover:text-rose-500 transition-colors">
          <ArrowRightOnRectangleIcon class="w-6 h-6" />
          <span class="text-[10px] font-semibold">Salir</span>
        </button>
      </nav>
      <!-- End mobile nav padding -->
      <div v-if="isPWA" class="h-24 md:h-0"></div>
    </main>

    <!-- ===== AppModal instances (self-contained — each has its own Teleport) ===== -->

    <!-- Shift selection modal -->
    <AppModal v-model="showShiftModal" title="¿Qué turno laboras hoy?" color="brand" :closable="true">
      <div class="flex flex-col gap-3 py-2">
        <button @click="confirmShift('DIURNO')" class="btn-primary btn-lg w-full"
          style="background-color: #f59e0b; border-color: #f59e0b; color: white;">
          DIURNO
        </button>
        <button @click="confirmShift('NOCTURNO')" class="btn-primary btn-lg w-full"
          style="background-color: #312e81; border-color: #312e81; color: white;">
          NOCTURNO
        </button>
      </div>
    </AppModal>

    <!-- Success modal -->
    <AppModal v-model="showSuccessModal" :title="checkType === 'entry' ? '¡Entrada Registrada!' : '¡Salida Registrada!'"
      :color="checkType === 'entry' ? 'success' : 'danger'" :closable="false" :closeOnBackdrop="false">
      <div class="flex flex-col items-center gap-3 py-2">
        <div class="w-14 h-14 rounded-2xl flex items-center justify-center"
          :class="checkType === 'entry' ? 'bg-emerald-100 dark:bg-emerald-500/15' : 'bg-rose-100 dark:bg-rose-500/15'">
          <CheckCircleIcon v-if="checkType === 'entry'" class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
          <HandRaisedIcon v-else class="w-8 h-8 text-rose-500  dark:text-rose-400" />
        </div>
        <div class="text-center">
          <p class="text-sm font-medium" style="color: var(--text-muted);">{{ registeredAt }}</p>
          <p class="text-xs mt-0.5" style="color: var(--text-muted);">Ubicación y fotografías guardadas exitosamente</p>
        </div>
      </div>
      <template #footer>
        <button @click="closeSuccess" class="btn-primary btn-lg w-full">Aceptar</button>
      </template>
    </AppModal>

    <!-- Logout confirmation modal -->
    <AppModal v-model="showLogoutModal" title="¿Cerrar sesión?" color="danger">
      <div class="flex flex-col items-center gap-4 py-2">
        <div class="w-12 h-12 rounded-full bg-rose-500/10 flex items-center justify-center">
          <ArrowRightOnRectangleIcon class="w-6 h-6 text-rose-500" />
        </div>
        <p class="text-sm text-center" style="color: var(--text-muted);">
          Tendrás que volver a ingresar tus credenciales para registrar asistencia.
        </p>
      </div>
      <template #footer>
        <button @click="showLogoutModal = false" class="btn-secondary flex-1">Cancelar</button>
        <button @click="confirmLogout" class="btn-danger flex-1">Salir</button>
      </template>
    </AppModal>

    <!-- ===== Modals ===== -->
    <Teleport to="body">

      <!-- Mobile Hamburger Menu Drawer -->
      <Transition name="slide-banner">
        <div v-if="showMobileMenu" class="fixed inset-0 z-50 flex justify-end">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showMobileMenu = false"></div>
          <div
            class="relative w-64 max-w-[80vw] h-full bg-white dark:bg-surface-950 shadow-2xl flex flex-col animate-slide-up sm:animate-in border-l border-slate-200 dark:border-surface-800">
            <div class="p-4 flex justify-between items-center border-b border-slate-100 dark:border-surface-800">
              <span class="font-bold text-slate-900 dark:text-white">Menú</span>
              <button @click="showMobileMenu = false"
                class="p-2 bg-slate-100 dark:bg-surface-800 rounded-full text-slate-500 dark:text-slate-400">
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <nav class="p-4 space-y-2 flex-1 overflow-y-auto">
              <router-link to="/checkin" @click="showMobileMenu = false"
                class="flex items-center gap-3 px-4 py-3 rounded-xl bg-brand-50 text-brand-600 dark:bg-brand-500/10 dark:text-brand-400 font-semibold transition-all">
                <MapPinIcon class="w-5 h-5" />
                Inicio
              </router-link>
              <router-link to="/history" @click="showMobileMenu = false"
                class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
                <ClockIcon class="w-5 h-5" />
                Historial
              </router-link>
              <router-link to="/qr" @click="showMobileMenu = false"
                class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 dark:text-slate-400 font-medium hover:bg-slate-50 dark:hover:bg-white/5 transition-all">
                <QrCodeIcon class="w-5 h-5" />
                Escanear QR
              </router-link>
            </nav>
            <div class="p-4 border-t border-slate-100 dark:border-surface-800 flex gap-2">
              <button @click="showLogoutModal = true; showMobileMenu = false"
                class="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-rose-500 dark:text-rose-400 font-semibold hover:bg-rose-50 dark:hover:bg-rose-500/10 transition-colors">
                <ArrowRightOnRectangleIcon class="w-5 h-5" />
                Cerrar Sesión
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Location error / guide modal -->
      <Transition name="modal">
        <div v-if="showLocationErrorModal"
          class="fixed inset-0 z-[60] flex items-end sm:items-center justify-center sm:p-4 bg-slate-900/60 dark:bg-slate-950/80 backdrop-blur-sm">
          <div
            class="w-full max-w-md glass-panel p-6 sm:p-8 animate-slide-up rounded-t-[2rem] sm:rounded-3xl border-b-0 sm:border-b border-t border-x border-white/20 dark:border-white/10 m-0">
            <div class="flex items-center justify-between mb-6">
              <div class="flex items-center gap-4">
                <div
                  class="w-12 h-12 rounded-2xl flex items-center justify-center flex-shrink-0 bg-rose-50 dark:bg-rose-500/10 border border-rose-100 dark:border-rose-500/20 text-rose-500">
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
              <div
                class="bg-brand-50 dark:bg-brand-500/5 border border-brand-100 dark:border-brand-500/20 rounded-2xl p-5 mb-6">
                <p class="text-sm font-bold mb-3 text-brand-600 dark:text-brand-400 uppercase tracking-wider">Cómo abrir
                  en Safari</p>
                <ol class="space-y-3">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span
                      class="w-6 h-6 rounded-full bg-brand-100 dark:bg-brand-500/20 text-brand-600 dark:text-brand-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">1</span>
                    <span class="mt-0.5">Toca el ícono de <strong class="text-slate-900 dark:text-white">tres puntos
                        (···)</strong> en la barra inferior.</span>
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span
                      class="w-6 h-6 rounded-full bg-brand-100 dark:bg-brand-500/20 text-brand-600 dark:text-brand-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">2</span>
                    <span class="mt-0.5">Selecciona <strong class="text-slate-900 dark:text-white">"Abrir en
                        navegador"</strong>.</span>
                  </li>
                </ol>
              </div>
            </template>

            <!-- GPS unavailable or timeout -->
            <template v-else-if="locationErrorType === 'unavailable' || locationErrorType === 'timeout'">
              <div
                class="bg-amber-50 dark:bg-amber-500/5 border border-amber-100 dark:border-amber-500/20 rounded-2xl p-5 mb-6">
                <ol class="space-y-3">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span
                      class="w-6 h-6 rounded-full bg-amber-100 dark:bg-amber-500/20 text-amber-600 dark:text-amber-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">1</span>
                    <span class="mt-0.5">Ve a un lugar con buena señal o al aire libre.</span>
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-3">
                    <span
                      class="w-6 h-6 rounded-full bg-amber-100 dark:bg-amber-500/20 text-amber-600 dark:text-amber-400 flex items-center justify-center flex-shrink-0 font-bold text-xs">2</span>
                    <span class="mt-0.5">Asegúrate de que el GPS esté <strong
                        class="text-slate-900 dark:text-white">activado</strong>.</span>
                  </li>
                </ol>
              </div>
            </template>

            <!-- Permission denied -->
            <template v-else>
              <!-- Paso 1: Configuración iOS global -->
              <div
                class="bg-brand-50 dark:bg-brand-500/5 border border-brand-100 dark:border-brand-500/20 rounded-2xl p-4 mb-4">
                <p class="text-xs font-bold mb-2.5 text-brand-600 dark:text-brand-400 uppercase tracking-wider">①
                  Config. del iPhone (Revisar primero)</p>
                <ol class="space-y-2">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-brand-500 font-bold mt-px">1.</span> Ve a <strong>Configuración → Privacidad →
                      Localización</strong>.
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-brand-500 font-bold mt-px">2.</span> Busca <strong>Safari</strong> y ponlo en
                    <strong>"Al usar la app"</strong>.
                  </li>
                </ol>
              </div>

              <!-- Paso 2: AA del sitio en Safari -->
              <div
                class="bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-2xl p-4 mb-6">
                <p class="text-xs font-bold mb-2.5 text-slate-500 dark:text-slate-400 uppercase tracking-wider">② En
                  esta página (Safari)</p>
                <ol class="space-y-2">
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-slate-400 dark:text-slate-500 font-bold mt-px">1.</span> Toca las letras
                    <strong>"AA"</strong> en la barra web.
                  </li>
                  <li class="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2">
                    <span class="text-slate-400 dark:text-slate-500 font-bold mt-px">2.</span> Ajustes del sitio web →
                    Ubicación <strong>"Permitir"</strong>.
                  </li>
                </ol>
              </div>
            </template>

            <div class="space-y-3 sm:space-y-0 sm:grid sm:grid-cols-2 sm:gap-3">
              <button @click="cancelProcess"
                class="btn-secondary btn-lg w-full order-last sm:order-first">Cancelar</button>
              <button v-if="locationErrorType === 'denied'" @click="reloadPage()" class="btn-primary btn-lg w-full">
                <ArrowPathRoundedSquareIcon class="w-5 h-5" /> Recargar
              </button>
              <button v-else-if="locationErrorType !== 'inapp'" @click="retryLocation()"
                class="btn-primary btn-lg w-full">
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
          <video ref="videoRef" autoplay playsinline muted class="absolute inset-0 w-full h-full object-cover"
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

      <!-- Success modal → moved outside Teleport as AppModal (has its own Teleport) -->

      <!-- Processing spinner -->
      <div v-if="processing && !showLocationErrorModal && !showCameraModal && !showSuccessModal"
        class="fixed inset-0 z-50 flex items-center justify-center"
        style="background: rgba(0,0,0,0.5); backdrop-filter: blur(4px);">
        <div class="glass-card p-6 text-center min-w-[200px]">
          <!-- GPS acquiring state -->
          <template v-if="gpsAcquiring">
            <div class="relative w-12 h-12 mx-auto mb-3">
              <div class="absolute inset-0 rounded-full border-2 border-amber-500/20 border-t-amber-500 animate-spin">
              </div>
              <div class="absolute inset-2 rounded-full border-2 border-amber-500/10 border-t-amber-400 animate-spin"
                style="animation-duration:1.5s;animation-direction:reverse"></div>
            </div>
            <p class="font-semibold text-amber-600 dark:text-amber-400">Obteniendo ubicación…</p>
            <p class="text-xs mt-1" style="color: var(--text-muted);">
              {{ gpsAccuracy != null ? `Precisión actual: ~${gpsAccuracy} m` : 'Esperando señal GPS…' }}
            </p>
            <p class="text-[10px] mt-1 text-slate-400">Buscando fix ≤ 30 m</p>
          </template>
          <!-- Generic processing state -->
          <template v-else>
            <div
              class="w-10 h-10 border-2 border-brand-500/30 border-t-brand-500 rounded-full animate-spin mx-auto mb-3">
            </div>
            <p class="font-semibold" style="color: var(--text);">Procesando...</p>
            <p class="text-xs mt-1" style="color: var(--text-muted);">Por favor espera</p>
          </template>
        </div>
      </div>

      <!-- Logout confirmation modal → moved outside Teleport as AppModal (has its own Teleport) -->

    </Teleport>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import { useToast } from '@/composables/useToast'
import AppModal from '@/components/AppModal.vue'
import {
  CheckCircleIcon, HandRaisedIcon, MapPinIcon, ArrowRightOnRectangleIcon,
  ExclamationTriangleIcon, XMarkIcon, LockClosedIcon, CloudArrowUpIcon,
  SignalIcon, MinusCircleIcon, CameraIcon, UserCircleIcon, BoltIcon, ArrowPathIcon,
  InformationCircleIcon, DevicePhoneMobileIcon, ArrowPathRoundedSquareIcon,
  QrCodeIcon, ClockIcon, Bars3BottomRightIcon
} from '@heroicons/vue/24/outline'

// ─── iziToast helpers ────────────────────────────────────────────────────────
const toast = useToast()

const auth = useAuthStore()
const router = useRouter()

const showLogoutModal = ref(false)
const isSyncing = ref(false)

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


// State — session state is ALWAYS derived from the server, never cached locally.
// localStorage was unreliable in PWA/iOS Safari (gets cleared after 7 days inactivity,
// private browsing, storage limits) causing the UI to show wrong state on reload.
const activeRecordId = ref(null)
const entryTime = ref(null)
const statusLoading = ref(true)
const statusError = ref(false)
const pendingCount = ref(0)   // offline-queued checks waiting to sync
const processing = ref(false)
const checkType = ref('entry')
const showLocationErrorModal = ref(false)
const locationErrorType = ref('denied')
const locationErrorCode = ref(0)
const showCameraModal = ref(false)
const showSuccessModal = ref(false)
const showShiftModal = ref(false)
const selectedShift = ref('')
const registeredAt = ref('')

const isPWA = ref(false)
const isOnline = ref(navigator.onLine)
function handleOnline() { isOnline.value = true; syncPendingChecks() }
function handleOffline() { isOnline.value = false }
const showMobileMenu = ref(false)

// Location
let locationWatchId = null
const locationPoints = []

// UI extras
const recoveredPoints = ref(0)
const showHiddenWarning = ref(false)
const isTracking = ref(false)
const gpsPermission = ref('prompt') // 'granted' | 'denied' | 'prompt'
const gpsAccuracy = ref(null)   // last known accuracy in meters (null = not yet acquired)
const gpsAcquiring = ref(false)  // true while waiting for a precise GPS fix

// GPS acquisition handles — stored at module level so onUnmounted can cancel them
// if the user navigates away while waiting for a GPS fix.
let pendingAcquireWatchId = null
let pendingAcquireTimeoutId = null

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

// ─── IndexedDB location buffer ────────────────────────────────────────────────
// IndexedDB is significantly more persistent than localStorage:
// - Survives iOS Safari PWA restarts (localStorage clears after 7-day inactivity)
// - Not affected by private-browsing storage limits
// - Larger storage quota
const IDB_NAME = 'paselista-db'
const IDB_VERSION = 2          // bumped: adds pending_checks store
const IDB_STORE = 'location_buffer'
const IDB_PENDING = 'pending_checks'

function openIDB() {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(IDB_NAME, IDB_VERSION)
    req.onupgradeneeded = (e) => {
      const db = e.target.result
      if (!db.objectStoreNames.contains(IDB_STORE)) {
        db.createObjectStore(IDB_STORE, { autoIncrement: true })
      }
      if (!db.objectStoreNames.contains(IDB_PENDING)) {
        db.createObjectStore(IDB_PENDING, { keyPath: 'id', autoIncrement: true })
      }
    }
    req.onsuccess = (e) => resolve(e.target.result)
    req.onerror = () => reject(req.error)
  })
}

// ─── GPS location buffer (active session) ────────────────────────────────────
async function saveLocationBuffer(points) {
  try {
    const db = await openIDB()
    const tx = db.transaction(IDB_STORE, 'readwrite')
    const store = tx.objectStore(IDB_STORE)
    store.clear()
    for (const p of points) store.add(p)
  } catch { /* silently fall back — non-critical */ }
}

async function loadLocationBuffer() {
  try {
    const db = await openIDB()
    return await new Promise((resolve) => {
      const req = db.transaction(IDB_STORE, 'readonly').objectStore(IDB_STORE).getAll()
      req.onsuccess = () => resolve(req.result ?? [])
      req.onerror = () => resolve([])
    })
  } catch { return [] }
}

async function clearLocationBuffer() {
  try {
    const db = await openIDB()
    db.transaction(IDB_STORE, 'readwrite').objectStore(IDB_STORE).clear()
  } catch { }
}

// ─── Offline pending checks queue ────────────────────────────────────────────
async function savePendingCheck(record) {
  try {
    const db = await openIDB()
    return await new Promise((resolve, reject) => {
      const req = db.transaction(IDB_PENDING, 'readwrite').objectStore(IDB_PENDING).add(record)
      req.onsuccess = () => resolve(req.result)
      req.onerror = () => reject(req.error)
    })
  } catch { return null }
}

async function loadPendingChecks() {
  try {
    const db = await openIDB()
    return await new Promise((resolve) => {
      const req = db.transaction(IDB_PENDING, 'readonly').objectStore(IDB_PENDING).getAll()
      req.onsuccess = () => resolve(req.result ?? [])
      req.onerror = () => resolve([])
    })
  } catch { return [] }
}

async function deletePendingCheck(id) {
  try {
    const db = await openIDB()
    db.transaction(IDB_PENDING, 'readwrite').objectStore(IDB_PENDING).delete(id)
  } catch { }
}
function handleVisibilityChange() {
  if (document.visibilityState === 'hidden') {
    saveLocationBuffer(locationPoints) // async fire-and-forget
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
  showShiftModal.value = true
}

function confirmShift(shift) {
  selectedShift.value = shift
  showShiftModal.value = false
  requestLocation()
}

function cancelProcess() {
  showLocationErrorModal.value = false
  processing.value = false
}

// requestLocation MUST start synchronously with a geolocation call so iOS
// Safari recognises it as originating from the user-gesture call stack.
// Strategy: use watchPosition to seek an accurate fix (≤ GOOD_ACCURACY_M).
// After MAX_WAIT_MS we accept the best reading obtained so far.
function requestLocation() {
  if (!navigator.geolocation) {
    locationErrorCode.value = 0
    locationErrorType.value = 'unavailable'
    showLocationErrorModal.value = true
    return
  }

  const GOOD_ACCURACY_M = 30   // accept immediately when ≤ 30 m
  const MAX_WAIT_MS = 20000 // fall back to best after 20 s
  const MAX_ACCURACY_M = 200  // reject readings worse than this

  let bestPosition = null
  let resolved = false
  let acquireWatchId = null
  let timeoutId = null

  async function finalize(position) {
    if (resolved) return
    resolved = true
    navigator.geolocation.clearWatch(acquireWatchId)
    clearTimeout(timeoutId)

    gpsAccuracy.value = Math.round(position.coords.accuracy)
    gpsAcquiring.value = false
    processing.value = false

    locationPoints.push({
      latitude: position.coords.latitude,
      longitude: position.coords.longitude,
      accuracy: position.coords.accuracy,
      recorded_at: new Date().toISOString()
    })

    // Start continuous background tracking, only store points ≤ MAX_ACCURACY_M
    locationWatchId = navigator.geolocation.watchPosition(
      pos => {
        gpsAccuracy.value = Math.round(pos.coords.accuracy)
        if (pos.coords.accuracy <= MAX_ACCURACY_M) {
          locationPoints.push({
            latitude: pos.coords.latitude,
            longitude: pos.coords.longitude,
            accuracy: pos.coords.accuracy,
            recorded_at: new Date().toISOString()
          })
          // Fix 3: cap array to avoid unbounded growth on long shifts (10h+ = thousands of pts)
          // Keep index 0 (entry point used by submitCheck) + most recent (MAX-1) points
          if (locationPoints.length > MAX_LOCATION_POINTS) {
            locationPoints.splice(1, locationPoints.length - MAX_LOCATION_POINTS)
          }
        }
      },
      null,
      { enableHighAccuracy: true, maximumAge: 0 }
    )
    isTracking.value = true
    requestWakeLock()

    cameraStep.value = 'site'
    showCameraModal.value = true
    await nextTick()
    await startCamera('environment')
  }

  function fail(err) {
    if (resolved) return
    resolved = true
    navigator.geolocation.clearWatch(acquireWatchId)
    clearTimeout(timeoutId)
    gpsAcquiring.value = false
    processing.value = false
    const code = err?.code ?? -1
    locationErrorCode.value = code
    if (code === 1) locationErrorType.value = 'denied'
    else if (code === 2) locationErrorType.value = 'unavailable'
    else locationErrorType.value = 'timeout'
    showLocationErrorModal.value = true
  }

  // ← watchPosition is the absolute first call — iOS Safari gesture chain preserved
  acquireWatchId = navigator.geolocation.watchPosition(
    (pos) => {
      gpsAccuracy.value = Math.round(pos.coords.accuracy)
      // Track best position seen so far
      if (!bestPosition || pos.coords.accuracy < bestPosition.coords.accuracy) {
        bestPosition = pos
      }
      // Good enough fix — proceed immediately
      if (pos.coords.accuracy <= GOOD_ACCURACY_M) {
        finalize(pos)
      }
    },
    fail,
    { enableHighAccuracy: true, maximumAge: 0, timeout: 30000 }
  )
  // Fix 4: expose to module scope so onUnmounted can cancel if user navigates away
  pendingAcquireWatchId = acquireWatchId

  // After MAX_WAIT_MS accept whatever best reading we have
  timeoutId = setTimeout(() => {
    if (!resolved) {
      if (bestPosition) finalize(bestPosition)
      else fail({ code: 3 }) // TIMEOUT
    }
  }, MAX_WAIT_MS)
  pendingAcquireTimeoutId = timeoutId

  processing.value = true
  gpsAcquiring.value = true
  gpsAccuracy.value = null
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

// Fix 3: cap in-memory GPS points to avoid unbounded growth on long shifts
const MAX_LOCATION_POINTS = 500

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

async function startCamera(facing = 'environment', { isSwitching = false } = {}) {
  stopCamera()
  facingFront.value = facing === 'user'
  try {
    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: facing, width: { ideal: 1280 }, height: { ideal: 720 } },
      audio: false
    })
    if (videoRef.value) videoRef.value.srcObject = stream
  } catch (err) {
    // Fix 2: when switching cameras keep modal open; when opening fresh, close it.
    if (!isSwitching) {
      showCameraModal.value = false
      stopLocationTracking()
      processing.value = false
      gpsAcquiring.value = false
    }
    if (err?.name === 'NotAllowedError' || err?.name === 'PermissionDeniedError') {
      toast.error('Permiso de cámara denegado. Habilítalo en la configuración del navegador.', 'Sin acceso a cámara')
    } else if (err?.name === 'NotFoundError') {
      toast.error('No se encontró una cámara en este dispositivo.', 'Sin cámara')
    } else if (isSwitching) {
      toast.warning('No se pudo cambiar de cámara.', 'Error de cámara')
    } else {
      toast.error('No se pudo acceder a la cámara. Intenta de nuevo.', 'Error de cámara')
    }
  }
}

function stopCamera() {
  if (stream) { stream.getTracks().forEach(t => t.stop()); stream = null }
}

async function switchCamera() {
  // Fix 2: pass isSwitching so modal stays open on failure
  await startCamera(facingFront.value ? 'environment' : 'user', { isSwitching: true })
}

function toggleFlash() {
  const track = stream?.getVideoTracks()[0]
  if (!track) return
  const caps = track.getCapabilities?.()
  if (caps?.torch) {
    flashOn.value = !flashOn.value
    track.applyConstraints({ advanced: [{ torch: flashOn.value }] })
  } else {
    toast.warning('Tu dispositivo no soporta flash desde el navegador.', 'Flash no disponible')
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

// Fix 3: re-encode a data URL at lower JPEG quality for IDB storage.
// Keeps original dimensions but reduces size from ~200 KB to ~50 KB.
function compressDataURL(dataURL, quality = 0.4) {
  if (!dataURL) return Promise.resolve(null)
  return new Promise(resolve => {
    const img = new Image()
    img.onload = () => {
      const canvas = document.createElement('canvas')
      canvas.width = img.width
      canvas.height = img.height
      canvas.getContext('2d').drawImage(img, 0, 0)
      resolve(canvas.toDataURL('image/jpeg', quality))
    }
    img.onerror = () => resolve(dataURL) // fallback: keep original on decode error
    img.src = dataURL
  })
}

async function submitCheck() {
  processing.value = true
  try {
    const now = new Date().toISOString()
    const formData = new FormData()
    formData.append('type', checkType.value)
    if (selectedShift.value) {
      formData.append('shift', selectedShift.value)
    }
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
    await clearLocationBuffer()
    stopLocationTracking()
    recoveredPoints.value = 0

    if (checkType.value === 'entry') {
      activeRecordId.value = recordId
      entryTime.value = new Date(now).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
    } else {
      activeRecordId.value = null
      entryTime.value = null
    }

    registeredAt.value = new Date(now).toLocaleString('es-MX')
    showSuccessModal.value = true
  } catch (err) {
    stopLocationTracking()
    const status = err?.response?.status
    const msg = err?.response?.data?.error || ''

    // Fix 1: detect captive-portal / any no-response scenario, not just navigator.onLine
    const isNetworkError = !navigator.onLine
      || err?.code === 'ERR_NETWORK'
      || err?.code === 'ECONNABORTED'
      || (err?.request && !err?.response) // request sent, zero response bytes received

    if (isNetworkError) {
      // ── Offline queue ──────────────────────────────────────────────────────
      // GPS and photos are already captured — save everything locally and
      // update the UI optimistically. Will sync automatically on reconnect.

      // Fix 3: enforce max queue depth to avoid filling IDB on low-end devices
      const MAX_PENDING = 5
      const existing = await loadPendingChecks()
      if (existing.length >= MAX_PENDING) {
        toast.error(`Cola offline llena (${MAX_PENDING} registros). Conéctate para sincronizar antes de continuar.`, 'Cola llena')
        return
      }

      // Fix 3: compress photos to ~40% quality before storing in IDB (~50 KB vs ~200 KB)
      const [compressedSite, compressedSelfie] = await Promise.all([
        compressDataURL(photoSite),
        compressDataURL(photoSelfie)
      ])

      const now2 = new Date().toISOString()
      await savePendingCheck({
        type: checkType.value,
        shift: selectedShift.value,
        timestamp: now2,
        latitude: locationPoints[0]?.latitude ?? 0,
        longitude: locationPoints[0]?.longitude ?? 0,
        photoSite: compressedSite,
        photoSelfie: compressedSelfie,
        queued_at: now2
      })
      pendingCount.value++

      // Update UI optimistically
      if (checkType.value === 'entry') {
        activeRecordId.value = 'pending'
        entryTime.value = new Date(now2).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
      } else {
        activeRecordId.value = null
        entryTime.value = null
      }
      registeredAt.value = new Date(now2).toLocaleString('es-MX')
      showSuccessModal.value = true   // show success modal — will sync later
      toast.warning(
        'Sin internet. El registro se guardó en tu dispositivo y se enviará automáticamente cuando haya conexión.',
        'Guardado offline'
      )
    } else if (status === 409) {
      // Conflict: active session mismatch — re-sync UI with server
      toast.warning(msg || 'Conflicto de sesión. La pantalla se ha actualizado.', 'Sesión actualizada')
      await syncStatus()
    } else if (status === 401) {
      // Handled by axios interceptor (redirect to login)
    } else {
      toast.error(msg || 'No se pudo registrar. Intenta de nuevo.', 'Error al registrar')
    }
  } finally {
    processing.value = false
    photoSite = null
    photoSelfie = null
  }
}

function closeSuccess() { showSuccessModal.value = false }

// ─── Lifecycle ────────────────────────────────────────────────────────────────
// Fetch active session status from server (single source of truth).
// On network failure: derives optimistic state from the IndexedDB pending queue
// so the user always sees the correct button (Entrada/Salida) even offline.
async function syncStatus() {
  statusLoading.value = true
  statusError.value = false
  try {
    const { data } = await api.get('/checks/status')
    if (data && data.active) {
      activeRecordId.value = data.record_id
      entryTime.value = new Date(data.entry_time).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
    } else {
      activeRecordId.value = null
      entryTime.value = null
    }
  } catch {
    // No server response — derive state from pending offline queue
    const pending = await loadPendingChecks()
    if (pending.length > 0) {
      // The last pending check determines the current optimistic state
      const last = pending[pending.length - 1]
      if (last.type === 'entry') {
        activeRecordId.value = 'pending'
        entryTime.value = new Date(last.timestamp).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
      } else {
        activeRecordId.value = null
        entryTime.value = null
      }
      // Show normal UI — pending banner already communicates the offline state
    } else {
      // No pending records and no server — truly unknown state
      statusError.value = true
    }
    if (!navigator.onLine) {
      toast.warning('Sin conexión. Mostrando estado guardado localmente.', 'Sin conexión')
    } else {
      toast.error('No se pudo conectar al servidor.', 'Error de conexión')
    }
  } finally {
    statusLoading.value = false
  }
}

function retryStatus() { syncStatus() }

// Attempt to submit all offline-queued checks to the server.
// Called automatically when the browser comes back online, and on mount.
async function syncPendingChecks() {
  const pending = await loadPendingChecks()
  if (pending.length === 0) return

  isSyncing.value = true
  toast.info('Sincronizando registros guardados offline…', 'Conexión restaurada')
  let synced = 0

  for (const p of pending) {
    try {
      const formData = new FormData()
      formData.append('type', p.type)
      if (p.shift) formData.append('shift', p.shift)
      formData.append('timestamp', p.timestamp)
      if (p.latitude) formData.append('latitude', String(p.latitude))
      if (p.longitude) formData.append('longitude', String(p.longitude))
      if (p.photoSite) formData.append('photo_site', dataURLtoBlob(p.photoSite), 'site.jpg')
      if (p.photoSelfie) formData.append('photo_selfie', dataURLtoBlob(p.photoSelfie), 'selfie.jpg')

      await api.post('/checks', formData)
      await deletePendingCheck(p.id)
      synced++
    } catch (err) {
      const status = err?.response?.status
      if (status === 409) {
        // Already registered (e.g. synced from another device) — discard silently
        await deletePendingCheck(p.id)
        synced++
      } else if (status === 401) {
        // Fix 1: JWT expired — stop immediately, records stay in IDB, user must re-login
        toast.warning(
          'Tu sesión expiró. Los registros pendientes se sincronizarán cuando vuelvas a iniciar sesión.',
          'Sesión expirada'
        )
        break // no point trying remaining records with same expired token
      }
      // Other errors (5xx, network): keep for next retry
    }
  }

  const failed = pending.length - synced
  if (synced > 0) {
    pendingCount.value = Math.max(0, pendingCount.value - synced)
    // Re-fetch real session state from server after sync
    await syncStatus()
    if (failed === 0) {
      toast.success(`${synced} registro${synced > 1 ? 's' : ''} sincronizado${synced > 1 ? 's' : ''} correctamente.`, '¡Sincronizado!')
    } else {
      // Fix 2: partial failure — some records synced, some still pending
      toast.warning(
        `${synced} sincronizado${synced > 1 ? 's' : ''}, ${failed} no pudo${failed > 1 ? 'ron' : ''} enviarse y se reintentará${failed > 1 ? 'n' : ''} más tarde.`,
        'Sincronización parcial'
      )
    }
  } else if (failed > 0) {
    // Fix 2: nothing synced — notify so the user knows it wasn't silent
    toast.warning('No se pudo sincronizar ningún registro. Se reintentará cuando haya conexión.', 'Sin sincronizar')
  }

  isSyncing.value = false
}

onMounted(async () => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  document.addEventListener('visibilitychange', handleVisibilityChange)

  // PWA standalone detection
  if (window.matchMedia('(display-mode: standalone)').matches || window.navigator.standalone) {
    isPWA.value = true
  }
  window.matchMedia('(display-mode: standalone)').addEventListener('change', (e) => {
    isPWA.value = e.matches
  })

  // Check GPS permission state (shows indicator even when not actively tracking)
  if (navigator.permissions) {
    try {
      const perm = await navigator.permissions.query({ name: 'geolocation' })
      gpsPermission.value = perm.state
      perm.onchange = () => { gpsPermission.value = perm.state }
    } catch { }
  }

  // Fetch session state from server — only source of truth
  await syncStatus()

  // Load pending offline count from IndexedDB
  const pending = await loadPendingChecks()
  pendingCount.value = pending.length

  // Recover buffered GPS points from IndexedDB if app was closed mid-session
  if (activeRecordId.value) {
    const buffered = await loadLocationBuffer()
    if (buffered.length > 0) {
      recoveredPoints.value = buffered.length
      locationPoints.push(...buffered)
    }
  }

  // Auto-sync pending offline checks when connection is restored
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)

  // If we're already online on mount and there are pending checks, sync now
  if (navigator.onLine && pending.length > 0) {
    await syncPendingChecks()
  }
})

onUnmounted(() => {
  clearInterval(timeInterval)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
  // Fix 4: cancel pending GPS acquisition if user navigated away before it resolved
  if (pendingAcquireWatchId !== null) {
    navigator.geolocation.clearWatch(pendingAcquireWatchId)
    pendingAcquireWatchId = null
  }
  if (pendingAcquireTimeoutId !== null) {
    clearTimeout(pendingAcquireTimeoutId)
    pendingAcquireTimeoutId = null
  }
  stopLocationTracking()
  stopCamera()
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.25s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.safe-top {
  padding-top: max(16px, env(safe-area-inset-top));
}

.safe-bottom {
  padding-bottom: max(24px, env(safe-area-inset-bottom));
}

.slide-banner-enter-active,
.slide-banner-leave-active {
  transition: all 0.3s ease;
}

.slide-banner-enter-from,
.slide-banner-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
