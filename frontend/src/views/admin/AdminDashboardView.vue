<template>
  <div class="admin-layout min-h-screen min-h-dvh" style="background-color: var(--bg);">

    <!-- Sidebar Desktop -->
    <aside class="sidebar hidden lg:flex flex-col fixed left-0 top-0 h-full z-30 w-64"
      style="background: var(--surface); border-right: 1px solid var(--border-subtle);">
      <!-- Logo -->
      <div class="px-6 py-6" style="border-bottom: 1px solid var(--border-subtle);">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center text-xl"
            style="background: linear-gradient(135deg, #6366f1, #8b5cf6);"><CheckCircleIcon class="w-5 h-5 text-white" /></div>
          <div>
            <p style="color: var(--text);" class="font-bold text-base leading-none">PaseLista</p>
            <p style="color: var(--text-muted);" class="text-xs mt-0.5">Panel Admin</p>
          </div>
        </div>
      </div>

      <!-- Nav -->
      <nav class="flex-1 px-3 py-4 space-y-1">
        <button v-for="item in navItems" :key="item.id"
          @click="activeTab = item.id"
          :class="activeTab === item.id ? 'text-white' : 'text-slate-400 hover:text-white hover:bg-white/5'"
          class="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all">
          <component :is="item.icon" class="w-5 h-5" />
          <span>{{ item.label }}</span>
          <div v-if="activeTab === item.id" class="ml-auto w-1.5 h-1.5 rounded-full bg-brand-400"></div>
        </button>
        <div v-if="activeTab !== 'users' && activeTab !== 'overview'" class="absolute left-64 top-0 bottom-0 w-0.5 bg-brand-500/20 pointer-events-none"></div>
      </nav>

      <!-- Admin info + logout -->
      <div class="px-4 py-4" style="border-top: 1px solid var(--border-subtle);">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center text-sm font-bold"
            style="background: linear-gradient(135deg, rgba(99,102,241,0.4), rgba(139,92,246,0.4)); border: 1px solid rgba(99,102,241,0.3);">
            {{ auth.user?.first_name?.[0] }}
          </div>
          <div class="min-w-0">
            <p class="text-sm font-medium truncate" style="color: var(--text);">{{ auth.user?.first_name }} {{ auth.user?.last_name }}</p>
            <p class="text-xs truncate" style="color: var(--text-muted);">Administrador</p>
          </div>
        </div>
        <div class="flex items-center justify-between mb-2">
          <span class="text-xs" style="color: var(--text-muted);">Tema</span>
          <ThemeToggle />
        </div>
        <button @click="handleLogout" class="btn-secondary btn-sm w-full text-rose-400 border-rose-500/20 hover:bg-rose-500/10">
          Cerrar Sesión
        </button>
      </div>
    </aside>

    <!-- Mobile top bar -->
    <header class="lg:hidden sticky top-0 z-20 px-4 py-3 flex items-center justify-between"
      style="background: var(--nav-bg); backdrop-filter: blur(12px); border-bottom: 1px solid var(--border-subtle);">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center"
          style="background: linear-gradient(135deg, #6366f1, #8b5cf6);"><CheckCircleIcon class="w-5 h-5 text-white" /></div>
        <span class="text-white font-bold text-sm">Admin · PaseLista</span>
      </div>
      <button @click="handleLogout" class="text-rose-400 text-sm font-medium">Salir</button>
    </header>

    <!-- Main content -->
    <main class="lg:pl-64 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-6">

        <!-- Page header -->
        <div class="mb-6 animate-in">
          <h1 class="text-2xl font-bold" style="color: var(--text);">
            {{ currentNavItem?.label }}
          </h1>
          <p class="text-sm mt-0.5" style="color: var(--text-muted);">{{ currentNavItem?.description }}</p>
        </div>

        <!-- ===== OVERVIEW TAB ===== -->
        <div v-if="activeTab === 'overview'" class="space-y-6 animate-in">
          <!-- Date filter -->
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-2">
              <label class="text-slate-400 text-sm">Fecha:</label>
              <input type="date" v-model="filters.date" @change="loadAll"
                class="input text-sm py-2 w-40" />
            </div>
            <button @click="setToday" class="badge-blue cursor-pointer hover:opacity-80 transition-opacity">Hoy</button>
            <button @click="clearDate" class="badge-gray cursor-pointer hover:opacity-80 transition-opacity">Todos</button>
          </div>

          <!-- Stats grid -->
          <div class="grid grid-cols-2 lg:grid-cols-5 gap-4">
            <div class="stat-card col-span-2 lg:col-span-1" style="background: rgba(99,102,241,0.1); border-color: rgba(99,102,241,0.25);">
              <p class="stat-label flex items-center gap-1"><UsersIcon class="w-4 h-4" /> Usuarios</p>
              <p class="stat-value text-brand-400">{{ stats.total_users }}</p>
            </div>
            <div class="stat-card" style="background: rgba(16,185,129,0.1); border-color: rgba(16,185,129,0.25);">
              <p class="stat-label flex items-center gap-1"><MapPinIcon class="w-4 h-4" /> Entradas</p>
              <p class="stat-value text-emerald-400">{{ stats.entries_total }}</p>
            </div>
            <div class="stat-card" style="background: rgba(244,63,94,0.1); border-color: rgba(244,63,94,0.25);">
              <p class="stat-label flex items-center gap-1"><ArrowRightOnRectangleIcon class="w-4 h-4" /> Salidas</p>
              <p class="stat-value text-rose-400">{{ stats.exits_total }}</p>
            </div>
            <div class="stat-card" style="background: rgba(251,191,36,0.1); border-color: rgba(251,191,36,0.25);">
              <p class="stat-label flex items-center gap-1"><SignalIcon class="w-4 h-4" /> Activos</p>
              <p class="stat-value text-amber-400">{{ stats.active_now }}</p>
            </div>
            <div class="stat-card" style="background: rgba(6,182,212,0.1); border-color: rgba(6,182,212,0.25);">
              <p class="stat-label flex items-center gap-1"><ClipboardDocumentListIcon class="w-4 h-4" /> Registros</p>
              <p class="stat-value text-cyan-400">{{ stats.records_today }}</p>
            </div>
          </div>

          <!-- Recent records -->
          <div class="glass-card overflow-hidden">
            <div class="px-5 py-4 flex items-center justify-between" style="border-bottom: 1px solid var(--border-subtle);">
              <h3 class="font-semibold" style="color: var(--text);">Registros recientes</h3>
              <button @click="activeTab = 'records'" class="text-brand-400 text-sm hover:text-brand-300 transition-colors">
                Ver todos →
              </button>
            </div>
            <div class="overflow-x-auto custom-scroll">
              <table class="data-table w-full">
                <thead>
                  <tr>
                    <th>Usuario</th>
                    <th class="hidden sm:table-cell">Proyecto</th>
                    <th>Tipo</th>
                    <th>Fecha / Hora</th>
                    <th class="hidden md:table-cell">GPS</th>
                    <th>Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-if="records.slice(0,8).length === 0">
                    <td colspan="6" class="text-center py-8" style="color: var(--text-muted);">Sin registros</td>
                  </tr>
                  <tr v-for="r in records.slice(0, 8)" :key="r.record_id"
                    :style="r.is_suspicious ? 'background: rgba(245,158,11,0.06);' : ''">
                    <td>
                      <div class="flex items-center gap-1.5">
                        <div>
                          <p class="font-medium text-sm" style="color: var(--text);">{{ r.first_name }} {{ r.last_name }}</p>
                          <p class="text-xs" style="color: var(--text-muted);">{{ r.email }}</p>
                        </div>
                        <ExclamationTriangleIcon v-if="r.is_suspicious"
                          class="w-4 h-4 text-amber-400 flex-shrink-0"
                          :title="r.suspicious_reason || 'Registro sospechoso'" />
                      </div>
                    </td>
                    <td class="hidden sm:table-cell">
                      <span class="text-xs line-clamp-1 max-w-[140px] block" style="color: var(--text-muted);">{{ r.project_name }}</span>
                    </td>
                    <td>
                      <span :class="r.type === 'entry' ? 'badge-green' : 'badge-red'" class="badge whitespace-nowrap">
                        {{ r.type === 'entry' ? 'Entrada' : 'Salida' }}
                      </span>
                    </td>
                    <td>
                      <p class="text-xs whitespace-nowrap" style="color: var(--text-muted);">{{ formatDate(r.timestamp) }}</p>
                    </td>
                    <td class="hidden md:table-cell">
                      <span :class="r.location_count > 0 ? 'badge-blue' : 'badge-gray'" class="badge">
                        {{ r.location_count }} pts
                      </span>
                    </td>
                    <td>
                      <div class="flex items-center gap-2">
                        <button v-if="r.has_site_photo || r.has_selfie_photo"
                          @click="openPhotosModal(r)"
                          class="text-xs px-2 py-1 rounded-lg transition-all flex-shrink-0"
                          style="color: #a78bfa;">Fotos</button>
                        <span v-else class="text-xs" style="color: var(--text-dim);">Sin fotos</span>
                        <button @click="openRouteModal(r)" :disabled="r.location_count === 0"
                          class="text-xs px-2 py-1 rounded-lg transition-all flex-shrink-0"
                          :style="r.location_count > 0 ? 'color: #818cf8;' : 'color: var(--text-dim); opacity: 0.4;'">Mapa</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- ===== RECORDS TAB ===== -->
        <div v-if="activeTab === 'records'" class="space-y-5 animate-in">
          <!-- Filter bar -->
          <div class="glass-card p-4">
            <div class="flex flex-wrap gap-3">
              <div class="flex items-center gap-2 flex-1 min-w-[150px]">
                <CalendarIcon class="w-4 h-4" style="color: var(--text-muted);" />
                <input type="date" v-model="filters.date" @change="loadRecords"
                  class="input text-sm py-2 flex-1" placeholder="Fecha" />
              </div>
              <div class="flex items-center gap-2 flex-1 min-w-[150px]">
                <BuildingOffice2Icon class="w-4 h-4" style="color: var(--text-muted);" />
                <select v-model="filters.project" @change="loadRecords"
                  class="input text-sm py-2 flex-1">
                  <option value="">Todos los proyectos</option>
                  <option v-for="p in projects" :key="p" :value="p">{{ p }}</option>
                </select>
              </div>
              <div class="flex items-center gap-2">
                <select v-model="filters.type" @change="loadRecords"
                  class="input text-sm py-2 w-36">
                  <option value="">Entradas y salidas</option>
                  <option value="entry">Solo entradas</option>
                  <option value="exit">Solo salidas</option>
                </select>
              </div>
              <div class="flex items-center gap-2 flex-1 min-w-[180px]">
                <MagnifyingGlassIcon class="w-4 h-4" style="color: var(--text-muted);" />
                <input v-model="filters.search" @input="debouncedLoad" type="text"
                  class="input text-sm py-2 flex-1" placeholder="Buscar por nombre o email..." />
              </div>
              <button @click="clearFilters" class="btn-secondary btn-sm">Limpiar</button>
            </div>
          </div>

          <!-- Count -->
          <div class="flex items-center justify-between">
            <p class="text-slate-400 text-sm">
              <span style="color: var(--text); font-weight: 600;">{{ records.length }}</span> <span style="color: var(--text-muted);">registros encontrados</span>
            </p>
            <div v-if="loadingRecords" class="flex items-center gap-2 text-slate-500 text-sm">
              <div class="w-3 h-3 border border-brand-500/40 border-t-brand-500 rounded-full animate-spin"></div>
              Cargando...
            </div>
          </div>

          <!-- Table -->
          <div class="glass-card overflow-hidden">
            <div v-if="records.length === 0 && !loadingRecords" class="text-center py-16">
              <ClipboardDocumentListIcon class="w-12 h-12 mx-auto mb-3 opacity-30" style="color: var(--text-muted);" />
              <p class="font-medium" style="color: var(--text-muted);">Sin resultados</p>
              <p class="text-sm mt-1" style="color: var(--text-dim);">Intenta cambiar los filtros</p>
            </div>
            <div v-else class="overflow-x-auto custom-scroll">
              <table class="data-table w-full">
                <thead>
                  <tr>
                    <th>Usuario</th>
                    <th class="hidden sm:table-cell">Proyecto</th>
                    <th>Tipo</th>
                    <th>Fecha / Hora</th>
                    <th class="hidden md:table-cell">GPS</th>
                    <th>Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="r in records" :key="r.record_id"
                    :style="r.is_suspicious ? 'background: rgba(245,158,11,0.06);' : ''">
                    <td>
                      <div class="flex items-center gap-1.5">
                        <div>
                          <p class="font-medium text-sm" style="color: var(--text);">{{ r.first_name }} {{ r.last_name }}</p>
                          <p class="text-xs" style="color: var(--text-muted);">{{ r.email }}</p>
                        </div>
                        <!-- Suspicious flag -->
                        <button v-if="r.is_suspicious" @click="openSuspiciousModal(r)"
                          class="flex-shrink-0 p-0.5 rounded"
                          title="Ver detalle de alerta">
                          <ExclamationTriangleIcon class="w-4 h-4 text-amber-400" />
                        </button>
                      </div>
                    </td>
                    <td class="hidden sm:table-cell">
                      <span class="text-xs line-clamp-1 max-w-[140px] block" style="color: var(--text-muted);">{{ r.project_name }}</span>
                    </td>
                    <td>
                      <span :class="r.type === 'entry' ? 'badge-green' : 'badge-red'" class="badge whitespace-nowrap">
                        {{ r.type === 'entry' ? 'Entrada' : 'Salida' }}
                      </span>
                    </td>
                    <td>
                      <p class="text-xs whitespace-nowrap" style="color: var(--text-muted);">{{ formatDate(r.timestamp) }}</p>
                    </td>
                    <td class="hidden md:table-cell">
                      <span :class="r.location_count > 0 ? 'badge-blue' : 'badge-gray'" class="badge">
                        {{ r.location_count }} pts
                      </span>
                    </td>
                    <td>
                      <div class="flex items-center gap-2">
                        <button v-if="r.has_site_photo || r.has_selfie_photo"
                          @click="openPhotosModal(r)"
                          class="text-xs px-2 py-1 rounded-lg transition-all flex-shrink-0"
                          style="color: #a78bfa;">Fotos</button>
                        <span v-else class="text-xs" style="color: var(--text-dim);">Sin fotos</span>
                        <button @click="openRouteModal(r)" :disabled="r.location_count === 0"
                          class="text-xs px-2 py-1 rounded-lg transition-all flex-shrink-0"
                          :style="r.location_count > 0 ? 'color: #818cf8;' : 'color: var(--text-dim); opacity: 0.4;'">Mapa</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- ===== USERS TAB ===== -->
        <div v-if="activeTab === 'users'" class="space-y-5 animate-in">
          <!-- Filter -->
          <div class="glass-card p-4">
            <div class="flex flex-wrap gap-3 items-center">
              <select v-model="filters.project" @change="loadUsers"
                class="input text-sm py-2 w-52">
                <option value="">Todos los proyectos</option>
                <option v-for="p in projects" :key="p" :value="p">{{ p }}</option>
              </select>
              <p class="text-slate-400 text-sm ml-auto">
                <span style="color: var(--text); font-weight: 600;">{{ users.length }}</span> <span style="color: var(--text-muted);">usuarios registrados</span>
              </p>
            </div>
          </div>

          <!-- Users grid -->
          <div v-if="users.length === 0" class="glass-card p-12 text-center">
            <UserCircleIcon class="w-12 h-12 mx-auto mb-3 opacity-30" style="color: var(--text-muted);" />
            <p style="color: var(--text-muted);">Sin usuarios registrados</p>
          </div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
            <div v-for="(user, i) in users" :key="user.id"
              class="glass-card p-5 animate-in"
              :style="`animation-delay: ${i * 0.03}s`">
              <div class="flex items-start gap-3 mb-4">
                <div class="w-12 h-12 rounded-2xl flex items-center justify-center text-lg font-bold flex-shrink-0"
                  style="background: linear-gradient(135deg, rgba(99,102,241,0.3), rgba(139,92,246,0.3)); border: 1px solid rgba(99,102,241,0.25);">
                  {{ user.first_name[0] }}{{ user.last_name[0] }}
                </div>
                <div class="min-w-0">
                  <p class="font-semibold truncate" style="color: var(--text);">{{ user.first_name }} {{ user.last_name }}</p>
                  <p class="text-xs truncate" style="color: var(--text-muted);">{{ user.email }}</p>
                </div>
              </div>
              <div class="space-y-2">
                <div class="flex items-center gap-2">
                  <span class="text-xs w-16" style="color: var(--text-dim);">Proyecto</span>
                  <span class="text-xs truncate" style="color: var(--text);">{{ user.project_name }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-16" style="color: var(--text-dim);">Registros</span>
                  <span class="badge-blue">{{ user.total_checks }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-16" style="color: var(--text-dim);">Registro</span>
                  <span class="text-xs" style="color: var(--text-muted);">{{ formatDateShort(user.created_at) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

      </div>
    </main>

    <!-- Mobile bottom nav -->
    <nav class="lg:hidden fixed bottom-0 left-0 right-0 z-20 px-2 py-2 safe-bottom flex justify-around"
      style="background: var(--nav-bg); backdrop-filter: blur(12px); border-top: 1px solid var(--border-subtle);">
      <button v-for="item in navItems" :key="item.id"
        @click="activeTab = item.id"
        :class="activeTab === item.id ? 'text-brand-400' : 'text-slate-500'"
        class="flex flex-col items-center gap-1 px-4 py-1 transition-colors">
        <component :is="item.icon" class="w-6 h-6" />
        <span class="text-xs font-medium">{{ item.label }}</span>
      </button>
    </nav>

    <!-- ===== MODALS ===== -->
    <Teleport to="body">

      <!-- Route modal -->
      <Transition name="modal">
        <div v-if="routeModal.show" class="fixed inset-0 z-50 flex items-center justify-center px-3 py-4"
          style="background: rgba(0,0,0,0.88); backdrop-filter: blur(10px);">
          <div class="w-full max-w-3xl glass-card flex flex-col animate-in overflow-hidden"
            style="background: var(--modal-bg); max-height: 94vh; border: 1px solid rgba(99,102,241,0.2);">

            <!-- Header premium -->
            <div class="flex items-center justify-between px-5 py-4 flex-shrink-0"
              style="background: linear-gradient(135deg, rgba(99,102,241,0.15), rgba(139,92,246,0.1)); border-bottom: 1px solid var(--border-subtle);">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                  style="background: linear-gradient(135deg, #6366f1, #8b5cf6);">
                  <MapPinIcon class="w-5 h-5 text-white" />
                </div>
                <div>
                  <h3 class="font-bold text-base" style="color: var(--text);">Recorrido GPS</h3>
                  <p class="text-xs mt-0.5" style="color: var(--text-muted);">{{ routeModal.user }}</p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <a v-if="routeModal.points.length > 0"
                  :href="googleMapsUrl(routeModal.points)"
                  target="_blank" rel="noopener"
                  class="btn-secondary btn-sm flex items-center gap-1.5 text-xs">
                  <MapPinIcon class="w-3.5 h-3.5" /> Ver en Google Maps
                </a>
                <button @click="routeModal.show = false; destroyMap()"
                  class="w-8 h-8 rounded-lg flex items-center justify-center transition-all hover:bg-white/10" style="color: var(--text-muted);">
                  <XMarkIcon class="w-5 h-5" />
                </button>
              </div>
            </div>

            <!-- Stats bar -->
            <div v-if="routeModal.points.length > 0"
              class="flex items-center gap-4 px-5 py-2.5 flex-shrink-0 flex-wrap"
              style="background: rgba(99,102,241,0.05); border-bottom: 1px solid var(--border-subtle);">
              <div class="flex items-center gap-1.5">
                <div class="w-3 h-3 rounded-full bg-emerald-400 flex-shrink-0"></div>
                <span class="text-xs" style="color: var(--text-muted);">Inicio:</span>
                <span class="text-xs font-semibold" style="color: var(--text);">{{ formatTimeOnly(routeModal.points[0]?.recorded_at) }}</span>
              </div>
              <div class="flex items-center gap-1.5">
                <div class="w-3 h-3 rounded-full bg-rose-400 flex-shrink-0"></div>
                <span class="text-xs" style="color: var(--text-muted);">Fin:</span>
                <span class="text-xs font-semibold" style="color: var(--text);">{{ formatTimeOnly(routeModal.points[routeModal.points.length-1]?.recorded_at) }}</span>
              </div>
              <div class="flex items-center gap-1.5">
                <SignalIcon class="w-3.5 h-3.5 text-brand-400" />
                <span class="text-xs font-semibold" style="color: var(--text);">{{ routeModal.points.length }} puntos GPS</span>
              </div>
              <div v-if="routeModal.duration" class="flex items-center gap-1.5">
                <span class="text-xs" style="color: var(--text-muted);">Duración:</span>
                <span class="text-xs font-semibold" style="color: var(--text);">{{ routeModal.duration }}</span>
              </div>
            </div>

            <!-- Map or empty -->
            <div v-if="routeModal.points.length === 0" class="text-center py-20 text-sm" style="color: var(--text-muted);">
              <MapPinIcon class="w-10 h-10 mx-auto mb-3 opacity-30" />
              Sin puntos GPS registrados
            </div>
            <div v-else class="flex-1 overflow-auto">
              <div ref="mapContainer" style="height: 420px; width: 100%;"></div>

              <!-- Address section -->
              <div class="px-5 py-4" style="border-top: 1px solid var(--border-subtle);">
                <div class="flex items-start gap-3">
                  <div class="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0 mt-0.5"
                    style="background: rgba(99,102,241,0.15); border: 1px solid rgba(99,102,241,0.25);">
                    <MapPinIcon class="w-4 h-4 text-brand-400" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-xs font-semibold mb-1" style="color: var(--text-muted);">📍 Dirección detectada (punto de inicio)</p>
                    <div v-if="routeModal.loadingAddress" class="flex items-center gap-2">
                      <div class="w-3 h-3 border border-brand-500/40 border-t-brand-500 rounded-full animate-spin"></div>
                      <span class="text-xs" style="color: var(--text-muted);">Obteniendo dirección...</span>
                    </div>
                    <p v-else-if="routeModal.address" class="text-sm font-medium" style="color: var(--text);">{{ routeModal.address }}</p>
                    <p v-else class="text-xs" style="color: var(--text-dim);">Dirección no disponible</p>
                    <p v-if="routeModal.points.length > 0" class="text-xs mt-1" style="color: var(--text-dim);">
                      Coordenadas: {{ routeModal.points[0].latitude.toFixed(6) }}, {{ routeModal.points[0].longitude.toFixed(6) }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

          </div>
        </div>
      </Transition>

      <!-- Suspicious alert modal -->
      <Transition name="modal">
        <div v-if="suspiciousModal.show" class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0,0,0,0.85); backdrop-filter: blur(8px);">
          <div class="w-full max-w-md glass-card animate-in" style="background: var(--modal-bg);">
            <div class="flex items-center justify-between px-5 py-4" style="border-bottom: 1px solid rgba(245,158,11,0.3);">
              <div class="flex items-center gap-2">
                <ExclamationTriangleIcon class="w-5 h-5 text-amber-400" />
                <h3 class="font-bold text-amber-400">Alerta de fraude detectado</h3>
              </div>
              <button @click="suspiciousModal.show = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center" style="color: var(--text-muted);">
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-5 space-y-4">
              <div class="rounded-xl p-4" style="background: rgba(245,158,11,0.08); border: 1px solid rgba(245,158,11,0.25);">
                <p class="text-sm font-semibold mb-1" style="color: var(--text);">{{ suspiciousModal.user }}</p>
                <p class="text-xs" style="color: var(--text-muted);">{{ suspiciousModal.date }}</p>
              </div>
              <div>
                <p class="text-xs font-semibold mb-2 text-amber-400">Razón de la alerta:</p>
                <p class="text-sm" style="color: var(--text);">{{ suspiciousModal.reason }}</p>
              </div>
              <div v-if="suspiciousModal.ipCountry" class="flex items-center gap-2">
                <span class="text-xs" style="color: var(--text-dim);">IP detectada en:</span>
                <span class="text-xs font-medium" style="color: var(--text);">{{ suspiciousModal.ipCity ? suspiciousModal.ipCity + ', ' : '' }}{{ suspiciousModal.ipCountry }}</span>
              </div>
              <div class="rounded-xl p-3" style="background: rgba(239,68,68,0.07); border: 1px solid rgba(239,68,68,0.2);">
                <p class="text-xs" style="color: var(--text-muted);">
                  El registro fue guardado normalmente. Esta alerta es solo informativa para que el administrador investigue.
                </p>
              </div>
              <button @click="suspiciousModal.show = false" class="btn-secondary btn-md w-full">Cerrar</button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Photos modal -->
      <Transition name="modal">
        <div v-if="photosModal.show" class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0,0,0,0.9); backdrop-filter: blur(8px);">
          <div class="w-full max-w-lg glass-card animate-in" style="background: var(--modal-bg);">
            <div class="flex items-center justify-between px-5 py-4" style="border-bottom: 1px solid var(--border-subtle);">
              <div>
                <h3 class="font-bold" style="color: var(--text);">Fotografías</h3>
                <p class="text-xs mt-0.5" style="color: var(--text-muted);">{{ photosModal.user }} · {{ photosModal.date }}</p>
              </div>
              <button @click="photosModal.show = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center transition-all" style="color: var(--text-muted);">
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-5 grid grid-cols-2 gap-4">
              <div>
                <p class="text-xs mb-2 flex items-center gap-1" style="color: var(--text-muted);"><BuildingOffice2Icon class="w-3 h-3" /> Fotografía del sitio</p>
                <div v-if="photosModal.site" class="rounded-xl overflow-hidden">
                  <img :src="photosModal.site" class="w-full aspect-square object-cover" />
                </div>
                <div v-else class="rounded-xl flex items-center justify-center aspect-square text-slate-600 text-sm"
                  style="background: var(--input-bg); border: 1px dashed var(--input-border);">
                  Sin foto
                </div>
              </div>
              <div>
                <p class="text-xs mb-2 flex items-center gap-1" style="color: var(--text-muted);"><UserCircleIcon class="w-3 h-3" /> Selfie</p>
                <div v-if="photosModal.selfie" class="rounded-xl overflow-hidden">
                  <img :src="photosModal.selfie" class="w-full aspect-square object-cover" />
                </div>
                <div v-else class="rounded-xl flex items-center justify-center aspect-square text-slate-600 text-sm"
                  style="background: var(--input-bg); border: 1px dashed var(--input-border);">
                  Sin foto
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import {
  CheckCircleIcon, MapPinIcon, ArrowRightOnRectangleIcon, SignalIcon,
  ClipboardDocumentListIcon, UsersIcon, ChartBarIcon, UserCircleIcon,
  BuildingOffice2Icon, CalendarIcon, MagnifyingGlassIcon, XMarkIcon, ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'
import { useThemeStore } from '@/stores/theme'
import ThemeToggle from '@/components/ThemeToggle.vue'

const auth = useAuthStore()
const theme = useThemeStore()
const router = useRouter()

const activeTab = ref('overview')
const navItems = [
  { id: 'overview', icon: ChartBarIcon, label: 'Resumen', description: 'Vista general del día' },
  { id: 'records', icon: ClipboardDocumentListIcon, label: 'Registros', description: 'Todas las entradas y salidas' },
  { id: 'users', icon: UsersIcon, label: 'Usuarios', description: 'Personas registradas en el sistema' }
]
const currentNavItem = computed(() => navItems.find(n => n.id === activeTab.value))

// Data
const stats = ref({ total_users: 0, records_today: 0, entries_total: 0, exits_total: 0, active_now: 0 })
const records = ref([])
const users = ref([])
const projects = ref([])
const loadingRecords = ref(false)

const filters = ref({ date: todayISO(), project: '', type: '', search: '' })

function todayISO() {
  return new Date().toISOString().split('T')[0]
}

function setToday() { filters.value.date = todayISO(); loadAll() }
function clearDate() { filters.value.date = ''; loadAll() }

function clearFilters() {
  filters.value = { date: todayISO(), project: '', type: '', search: '' }
  loadAll()
}

let searchTimer = null
function debouncedLoad() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(loadRecords, 400)
}

async function loadAll() {
  await Promise.all([loadStats(), loadRecords(), loadUsers(), loadProjects()])
}

async function loadStats() {
  try {
    const params = filters.value.date ? { date: filters.value.date } : {}
    const { data } = await api.get('/admin/stats', { params })
    stats.value = data
  } catch {}
}

async function loadRecords() {
  loadingRecords.value = true
  try {
    const params = {}
    if (filters.value.date) params.date = filters.value.date
    if (filters.value.project) params.project = filters.value.project
    if (filters.value.type) params.type = filters.value.type
    if (filters.value.search) params.search = filters.value.search
    const { data } = await api.get('/admin/records', { params })
    records.value = data
  } catch {} finally {
    loadingRecords.value = false
  }
}

async function loadUsers() {
  try {
    const params = filters.value.project ? { project: filters.value.project } : {}
    const { data } = await api.get('/admin/users', { params })
    users.value = data
  } catch {}
}

async function loadProjects() {
  try {
    const { data } = await api.get('/admin/projects')
    projects.value = data
  } catch {}
}

// Route modal + Leaflet map
const mapContainer = ref(null)
let leafletMap = null
const routeModal = ref({ show: false, user: '', points: [], duration: '', address: '', loadingAddress: false })

function googleMapsUrl(points) {
  if (!points.length) return '#'
  if (points.length === 1) {
    return `https://www.google.com/maps?q=${points[0].latitude},${points[0].longitude}`
  }
  const origin = `${points[0].latitude},${points[0].longitude}`
  const dest = `${points[points.length - 1].latitude},${points[points.length - 1].longitude}`
  const waypoints = points.slice(1, -1).map(p => `${p.latitude},${p.longitude}`).join('|')
  const base = `https://www.google.com/maps/dir/?api=1&origin=${origin}&destination=${dest}&travelmode=walking`
  return waypoints ? `${base}&waypoints=${encodeURIComponent(waypoints)}` : base
}

function destroyMap() {
  if (leafletMap) { leafletMap.remove(); leafletMap = null }
}

watch(() => routeModal.value.show, async (show) => {
  if (!show) { destroyMap(); return }
  if (!routeModal.value.points.length) return
  await nextTick()
  if (!mapContainer.value) return
  destroyMap()
  const points = routeModal.value.points
  const coords = points.map(p => [p.latitude, p.longitude])

  // Dynamic tile layer based on theme
  leafletMap = L.map(mapContainer.value, { zoomControl: true }).setView(coords[0], 17)
  
  const tileUrl = theme.isDark 
    ? 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
    : 'https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png'

  L.tileLayer(tileUrl, {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/">CARTO</a>',
    subdomains: 'abcd',
    maxZoom: 20
  }).addTo(leafletMap)

  // Route line — glowing effect
  L.polyline(coords, { color: '#6366f1', weight: 5, opacity: 0.9 }).addTo(leafletMap)
  L.polyline(coords, { color: '#a78bfa', weight: 2, opacity: 0.4 }).addTo(leafletMap)

  // Start marker (green)
  const startIcon = L.divIcon({
    html: `<div style="width:18px;height:18px;background:#10b981;border:3px solid white;border-radius:50%;box-shadow:0 0 8px rgba(16,185,129,0.8);"></div>`,
    className: '', iconAnchor: [9, 9]
  })
  L.marker(coords[0], { icon: startIcon })
    .bindPopup(`<b style="color:#10b981">🟢 Inicio</b><br><span style="font-size:11px">${formatDate(points[0].recorded_at)}</span>`)
    .addTo(leafletMap)

  // End marker (red)
  if (coords.length > 1) {
    const endIcon = L.divIcon({
      html: `<div style="width:18px;height:18px;background:#f43f5e;border:3px solid white;border-radius:50%;box-shadow:0 0 8px rgba(244,63,94,0.8);"></div>`,
      className: '', iconAnchor: [9, 9]
    })
    L.marker(coords[coords.length - 1], { icon: endIcon })
      .bindPopup(`<b style="color:#f43f5e">🔴 Fin</b><br><span style="font-size:11px">${formatDate(points[points.length - 1].recorded_at)}</span>`)
      .addTo(leafletMap)
  }

  leafletMap.fitBounds(L.polyline(coords).getBounds(), { padding: [32, 32] })

  // Reverse geocoding with Nominatim
  routeModal.value.loadingAddress = true
  routeModal.value.address = ''
  try {
    const { latitude: lat, longitude: lon } = points[0]
    const res = await fetch(
      `https://nominatim.openstreetmap.org/reverse?lat=${lat}&lon=${lon}&format=json&addressdetails=1`,
      { headers: { 'Accept-Language': 'es' } }
    )
    const geo = await res.json()
    routeModal.value.address = geo.display_name || ''
  } catch {
    routeModal.value.address = ''
  } finally {
    routeModal.value.loadingAddress = false
  }
})

async function openRouteModal(record) {
  try {
    const { data } = await api.get(`/admin/records/${record.record_id}/route`)
    // Calculate duration
    let duration = ''
    if (data.length >= 2) {
      const ms = new Date(data[data.length - 1].recorded_at) - new Date(data[0].recorded_at)
      const totalMin = Math.round(ms / 60000)
      duration = totalMin >= 60
        ? `${Math.floor(totalMin / 60)}h ${totalMin % 60}min`
        : `${totalMin} min`
    }
    routeModal.value = {
      show: true,
      user: `${record.first_name} ${record.last_name}`,
      points: data,
      duration,
      address: '',
      loadingAddress: false
    }
  } catch {}
}

// Suspicious modal
const suspiciousModal = ref({ show: false, user: '', date: '', reason: '', ipCountry: '', ipCity: '' })
function openSuspiciousModal(record) {
  suspiciousModal.value = {
    show: true,
    user: `${record.first_name} ${record.last_name}`,
    date: formatDate(record.timestamp),
    reason: record.suspicious_reason || 'Sin detalle disponible',
    ipCountry: record.ip_country || '',
    ipCity: record.ip_city || ''
  }
}

// Photos modal
const photosModal = ref({ show: false, user: '', date: '', site: '', selfie: '' })
async function openPhotosModal(record) {
  photosModal.value = {
    show: true,
    user: `${record.first_name} ${record.last_name}`,
    date: formatDate(record.timestamp),
    site: '',
    selfie: ''
  }
  try {
    const { data } = await api.get(`/admin/records/${record.record_id}/photos`)
    photosModal.value.site = data.photo_site_path || ''
    photosModal.value.selfie = data.photo_selfie_path || ''
  } catch {}
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('es-MX', {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

function formatDateShort(iso) {
  return new Date(iso).toLocaleDateString('es-MX', { year: 'numeric', month: 'short', day: 'numeric' })
}

function formatTimeOnly(iso) {
  if (!iso) return '--:--'
  return new Date(iso).toLocaleTimeString('es-MX', { hour: '2-digit', minute: '2-digit' })
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}

onMounted(() => loadAll())
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.25s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.safe-bottom { padding-bottom: max(8px, env(safe-area-inset-bottom)); }
.line-clamp-1 { display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; }

@media (min-width: 1024px) {
  .admin-layout { display: block; }
}
</style>
