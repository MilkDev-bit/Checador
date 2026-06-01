<template>
  <div
    class="admin-layout min-h-screen min-h-dvh"
    style="background-color: var(--bg)"
  >
    <!-- Sidebar Desktop -->
    <aside
      class="sidebar hidden lg:flex flex-col fixed left-0 top-0 h-full z-30 w-64"
      style="
        background: var(--surface);
        border-right: 1px solid var(--border-subtle);
      "
    >
      <!-- Logo -->
      <div
        class="px-6 py-6"
        style="border-bottom: 1px solid var(--border-subtle)"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-xl overflow-hidden flex items-center justify-center text-xl"
          >
            <img
              src="/checador-logo.png"
              alt="PaseLista Logo"
              class="w-full h-full object-cover"
            />
          </div>
          <div></div>
        </div>
      </div>

      <!-- Nav -->
      <nav class="flex-1 px-3 py-4 space-y-1">
        <button
          v-for="item in navItems"
          :key="item.id"
          @click="activeTab = item.id"
          :class="
            activeTab === item.id
              ? 'bg-gradient-to-r from-brand-500/90 to-violet-600/90 text-white shadow-md'
              : 'text-[var(--text-muted)] hover:text-[var(--text)] hover:bg-[var(--input-bg)]'
          "
          class="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-semibold transition-all"
        >
          <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
          <span>{{ item.label }}</span>
          <div
            v-if="activeTab === item.id"
            class="ml-auto w-1.5 h-1.5 rounded-full bg-white/70"
          ></div>
        </button>
        <div
          v-if="activeTab !== 'users' && activeTab !== 'overview'"
          class="absolute left-64 top-0 bottom-0 w-0.5 bg-brand-500/20 pointer-events-none"
        ></div>
      </nav>

      <!-- Admin info + logout -->
      <div class="px-4 py-4" style="border-top: 1px solid var(--border-subtle)">
        <div class="flex items-center gap-3 mb-3">
          <div
            class="w-9 h-9 rounded-xl flex items-center justify-center text-sm font-bold overflow-hidden"
            style="
              background: linear-gradient(
                135deg,
                rgba(99, 102, 241, 0.4),
                rgba(139, 92, 246, 0.4)
              );
              border: 1px solid rgba(99, 102, 241, 0.3);
            "
          >
            <img
              v-if="auth.user?.avatar_url"
              :src="auth.user.avatar_url"
              loading="lazy"
              class="w-full h-full object-cover"
              alt="Perfil"
            />
            <span v-else>{{ auth.user?.first_name?.[0] }}</span>
          </div>
          <div class="min-w-0">
            <p class="text-sm font-medium truncate" style="color: var(--text)">
              {{ auth.user?.first_name }} {{ auth.user?.last_name }}
            </p>
            <p class="text-xs truncate" style="color: var(--text-muted)">
              Administrador
            </p>
          </div>
        </div>
        <button
          @click="showLogoutModal = true"
          class="btn-secondary btn-sm w-full text-rose-400 border-rose-500/20 hover:bg-rose-500/10"
        >
          Cerrar Sesión
        </button>
      </div>
    </aside>

    <!-- Mobile top bar -->
    <header
      class="lg:hidden sticky top-0 z-20 px-4 py-3 safe-top flex items-center justify-between"
      style="
        background: var(--nav-bg);
        backdrop-filter: blur(12px);
        border-bottom: 1px solid var(--border-subtle);
      "
    >
      <div class="flex items-center gap-2">
        <div
          class="w-8 h-8 rounded-lg overflow-hidden flex items-center justify-center"
        >
          <img
            src="/checador-logo.png"
            alt="Logo"
            class="w-full h-full object-cover"
          />
        </div>
        <span class="text-white font-bold text-sm">Admin · PaseLista</span>
      </div>
      <button
        @click="showLogoutModal = true"
        class="text-rose-400 text-sm font-medium"
      >
        Salir
      </button>
    </header>

    <!-- Main content -->
    <main class="lg:pl-64 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 pt-6 pb-32 lg:pb-6">
        <!-- Page header -->
        <div class="mb-6 animate-in">
          <h1 class="text-2xl font-bold" style="color: var(--text)">
            {{ currentNavItem?.label }}
          </h1>
          <p class="text-sm mt-0.5" style="color: var(--text-muted)">
            {{ currentNavItem?.description }}
          </p>
        </div>

        <!-- ===== OVERVIEW TAB ===== -->
        <div v-if="activeTab === 'overview'" class="space-y-6 animate-in">
          <!-- Date filter -->
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-2">
              <label class="text-slate-400 text-sm">Fecha:</label>
              <input
                type="date"
                v-model="filters.date_from"
                @change="loadAll"
                class="input text-sm py-2 w-40"
              />
            </div>
            <button
              @click="setToday"
              class="badge-blue cursor-pointer hover:opacity-80 transition-opacity"
            >
              Hoy
            </button>
            <button
              @click="clearDate"
              class="badge-gray cursor-pointer hover:opacity-80 transition-opacity"
            >
              Todos
            </button>
          </div>

          <!-- Stats grid -->
          <div class="grid grid-cols-2 lg:grid-cols-5 gap-4">
            <div
              class="stat-card col-span-2 lg:col-span-1"
              style="
                background: rgba(99, 102, 241, 0.1);
                border-color: rgba(99, 102, 241, 0.3);
              "
            >
              <p class="stat-label flex items-center gap-1">
                <UsersIcon class="w-4 h-4 text-brand-500 dark:text-brand-400" />
                Usuarios
              </p>
              <p class="stat-value text-brand-600 dark:text-brand-400">
                {{ stats.total_users }}
              </p>
            </div>
            <div
              class="stat-card"
              style="
                background: rgba(16, 185, 129, 0.1);
                border-color: rgba(16, 185, 129, 0.3);
              "
            >
              <p class="stat-label flex items-center gap-1">
                <MapPinIcon
                  class="w-4 h-4 text-emerald-600 dark:text-emerald-400"
                />
                Entradas
              </p>
              <p class="stat-value text-emerald-600 dark:text-emerald-400">
                {{ stats.entries_total }}
              </p>
            </div>
            <div
              class="stat-card"
              style="
                background: rgba(244, 63, 94, 0.1);
                border-color: rgba(244, 63, 94, 0.3);
              "
            >
              <p class="stat-label flex items-center gap-1">
                <ArrowRightOnRectangleIcon
                  class="w-4 h-4 text-rose-600 dark:text-rose-400"
                />
                Salidas
              </p>
              <p class="stat-value text-rose-600 dark:text-rose-400">
                {{ stats.exits_total }}
              </p>
            </div>
            <div
              class="stat-card"
              style="
                background: rgba(251, 191, 36, 0.1);
                border-color: rgba(251, 191, 36, 0.3);
              "
            >
              <p class="stat-label flex items-center gap-1">
                <SignalIcon
                  class="w-4 h-4 text-amber-600 dark:text-amber-400"
                />
                Activos
              </p>
              <p class="stat-value text-amber-600 dark:text-amber-400">
                {{ stats.active_now }}
              </p>
            </div>
            <div
              class="stat-card"
              style="
                background: rgba(6, 182, 212, 0.1);
                border-color: rgba(6, 182, 212, 0.3);
              "
            >
              <p class="stat-label flex items-center gap-1">
                <ClipboardDocumentListIcon
                  class="w-4 h-4 text-cyan-600 dark:text-cyan-400"
                />
                Registros
              </p>
              <p class="stat-value text-cyan-600 dark:text-cyan-400">
                {{ stats.records_today }}
              </p>
            </div>
          </div>

          <!-- Recent records -->
          <div class="glass-card overflow-hidden">
            <div
              class="px-5 py-4 flex items-center justify-between"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <h3 class="font-semibold" style="color: var(--text)">
                Registros recientes
              </h3>
              <button
                @click="activeTab = 'records'"
                class="text-brand-600 dark:text-brand-400 text-sm hover:text-brand-700 dark:hover:text-brand-300 font-semibold transition-colors"
              >
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
                  <tr v-if="records.slice(0, 8).length === 0">
                    <td
                      colspan="6"
                      class="text-center py-8"
                      style="color: var(--text-muted)"
                    >
                      Sin registros
                    </td>
                  </tr>
                  <tr
                    v-for="r in records.slice(0, 8)"
                    :key="r.record_id"
                    :style="
                      r.is_suspicious
                        ? 'background: rgba(245,158,11,0.06);'
                        : ''
                    "
                  >
                    <td>
                      <div class="flex items-center gap-1.5">
                        <div>
                          <p
                            class="font-medium text-sm"
                            style="color: var(--text)"
                          >
                            {{ toTitleCase(r.first_name) }}
                            {{ toTitleCase(r.last_name) }}
                          </p>
                          <p class="text-xs" style="color: var(--text-muted)">
                            {{ r.email }}
                          </p>
                        </div>
                        <ExclamationTriangleIcon
                          v-if="r.is_suspicious"
                          class="w-4 h-4 text-amber-400 flex-shrink-0"
                          :title="r.suspicious_reason || 'Registro sospechoso'"
                        />
                      </div>
                    </td>
                    <td class="hidden sm:table-cell">
                      <span
                        class="text-xs line-clamp-1 max-w-[140px] block"
                        style="color: var(--text-muted)"
                        >{{ r.project_name }}</span
                      >
                    </td>
                    <td>
                      <span
                        :class="
                          r.type === 'entry' ? 'badge-green' : 'badge-red'
                        "
                        class="badge whitespace-nowrap"
                      >
                        {{ r.type === "entry" ? "Entrada" : "Salida" }}
                      </span>
                    </td>
                    <td>
                      <p
                        class="text-xs whitespace-nowrap"
                        style="color: var(--text-muted)"
                      >
                        {{ formatDate(r.timestamp) }}
                      </p>
                    </td>
                    <td class="hidden md:table-cell">
                      <span
                        :class="
                          r.location_count > 0 ? 'badge-blue' : 'badge-gray'
                        "
                        class="badge"
                      >
                        {{ r.location_count }} pts
                      </span>
                    </td>
                    <td>
                      <div class="flex items-center gap-2">
                        <button
                          v-if="r.has_site_photo || r.has_selfie_photo"
                          @click="openPhotosModal(r)"
                          class="text-xs px-2 py-1 rounded-lg font-semibold transition-all flex-shrink-0 bg-violet-100 text-violet-700 hover:bg-violet-200 dark:bg-violet-500/15 dark:text-violet-400"
                        >
                          Fotos
                        </button>
                        <span
                          v-else
                          class="text-xs"
                          style="color: var(--text-dim)"
                          >Sin fotos</span
                        >
                        <button
                          @click="openRouteModal(r)"
                          :disabled="r.location_count === 0"
                          class="text-xs px-2 py-1 rounded-lg font-semibold transition-all flex-shrink-0 bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-500/15 dark:text-indigo-400 disabled:opacity-40 disabled:cursor-not-allowed"
                        >
                          Mapa
                        </button>
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
          <div class="glass-card p-4 space-y-3">
            <!-- Fila 1: Rango de fechas + accesos rápidos -->
            <div class="flex flex-wrap items-center gap-3">
              <div class="flex items-center gap-2">
                <CalendarIcon
                  class="w-4 h-4 flex-shrink-0"
                  style="color: var(--text-muted)"
                />
                <span
                  class="text-xs font-medium"
                  style="color: var(--text-muted)"
                  >Desde</span
                >
                <input
                  type="date"
                  v-model="filters.date_from"
                  @change="loadRecords"
                  class="input text-sm py-2 w-40"
                />
              </div>
              <div class="flex items-center gap-2">
                <span
                  class="text-xs font-medium"
                  style="color: var(--text-muted)"
                  >Hasta</span
                >
                <input
                  type="date"
                  v-model="filters.date_to"
                  @change="loadRecords"
                  class="input text-sm py-2 w-40"
                />
              </div>
              <button
                @click="setToday"
                class="badge-blue cursor-pointer hover:opacity-80 transition-opacity text-xs px-3 py-1.5 rounded-lg"
              >
                Hoy
              </button>
              <button
                @click="clearDate"
                class="badge-gray cursor-pointer hover:opacity-80 transition-opacity text-xs px-3 py-1.5 rounded-lg"
              >
                Todos
              </button>
            </div>
            <!-- Fila 2: Proyecto, tipo, búsqueda, limpiar -->
            <div class="flex flex-wrap items-center gap-3">
              <div class="flex items-center gap-2 flex-1 min-w-[150px]">
                <BuildingOffice2Icon
                  class="w-4 h-4 flex-shrink-0"
                  style="color: var(--text-muted)"
                />
                <select
                  v-model="filters.project"
                  @change="loadRecords"
                  class="input text-sm py-2 flex-1"
                >
                  <option value="">Todos los proyectos</option>
                  <option v-for="p in projects" :key="p" :value="p">
                    {{ p }}
                  </option>
                </select>
              </div>
              <div class="flex items-center gap-2">
                <select
                  v-model="filters.type"
                  @change="loadRecords"
                  class="input text-sm py-2 w-44"
                >
                  <option value="">Entradas y salidas</option>
                  <option value="entry">Solo entradas</option>
                  <option value="exit">Solo salidas</option>
                </select>
              </div>
              <div class="flex items-center gap-2 flex-1 min-w-[150px]">
                <UserCircleIcon
                  class="w-4 h-4 flex-shrink-0"
                  style="color: var(--text-muted)"
                />
                <input
                  v-model="filters.name"
                  @input="debouncedLoad"
                  type="text"
                  class="input text-sm py-2 flex-1"
                  placeholder="Buscar por nombre..."
                />
              </div>
              <div class="flex items-center gap-2 flex-1 min-w-[180px]">
                <MagnifyingGlassIcon
                  class="w-4 h-4 flex-shrink-0"
                  style="color: var(--text-muted)"
                />
                <input
                  v-model="filters.search"
                  @input="debouncedLoad"
                  type="text"
                  class="input text-sm py-2 flex-1"
                  placeholder="Buscar por email..."
                />
              </div>
              <button @click="clearFilters" class="btn-secondary btn-sm">
                Limpiar
              </button>
            </div>
          </div>

          <!-- Count + Excel -->
          <div class="flex items-center justify-between gap-3 flex-wrap">
            <p class="text-slate-400 text-sm">
              <span style="color: var(--text); font-weight: 600">{{
                records.length
              }}</span>
              <span style="color: var(--text-muted)"
                >registros encontrados</span
              >
            </p>
            <div class="flex items-center gap-3">
              <div
                v-if="loadingRecords"
                class="flex items-center gap-2 text-slate-500 text-sm"
              >
                <div
                  class="w-3 h-3 border border-brand-500/40 border-t-brand-500 rounded-full animate-spin"
                ></div>
                Cargando...
              </div>
              <button
                v-if="records.length > 0"
                @click="exportExcel"
                class="btn-secondary btn-sm flex items-center gap-1.5 text-green-400 border-green-500/20 hover:bg-green-500/10"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5 5 5-5M12 3v12"
                  />
                </svg>
                Descargar Excel
              </button>
            </div>
          </div>

          <!-- Table -->
          <div class="glass-card overflow-hidden">
            <div
              v-if="records.length === 0 && !loadingRecords"
              class="text-center py-16"
            >
              <ClipboardDocumentListIcon
                class="w-12 h-12 mx-auto mb-3 opacity-30"
                style="color: var(--text-muted)"
              />
              <p class="font-medium" style="color: var(--text-muted)">
                Sin resultados
              </p>
              <p class="text-sm mt-1" style="color: var(--text-dim)">
                Intenta cambiar los filtros
              </p>
            </div>
            <div v-else class="overflow-x-auto custom-scroll">
              <table class="data-table w-full">
                <thead>
                  <tr>
                    <th>Usuario</th>
                    <th class="hidden sm:table-cell">Proyecto</th>
                    <th>Tipo</th>
                    <th>Fecha / Hora</th>
                    <th class="hidden lg:table-cell">Horario</th>
                    <th class="hidden lg:table-cell">Retraso / Adelanto</th>
                    <th class="hidden md:table-cell">GPS</th>
                    <th class="hidden xl:table-cell">Comentario del día</th>
                    <th>Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="r in paginatedRecords"
                    :key="r.record_id"
                    :style="
                      r.is_suspicious
                        ? 'background: rgba(245,158,11,0.06);'
                        : ''
                    "
                  >
                    <td>
                      <div class="flex items-center gap-1.5">
                        <div>
                          <p
                            class="font-medium text-sm"
                            style="color: var(--text)"
                          >
                            {{ toTitleCase(r.first_name) }}
                            {{ toTitleCase(r.last_name) }}
                          </p>
                          <p class="text-xs" style="color: var(--text-muted)">
                            {{ r.email }}
                          </p>
                        </div>
                        <!-- Suspicious flag -->
                        <button
                          v-if="r.is_suspicious"
                          @click="openSuspiciousModal(r)"
                          class="flex-shrink-0 p-0.5 rounded"
                          title="Ver detalle de alerta"
                        >
                          <ExclamationTriangleIcon
                            class="w-4 h-4 text-amber-400"
                          />
                        </button>
                      </div>
                    </td>
                    <td class="hidden sm:table-cell">
                      <span
                        class="text-xs line-clamp-1 max-w-[140px] block"
                        style="color: var(--text-muted)"
                        >{{ r.project_name }}</span
                      >
                    </td>
                    <td>
                      <span
                        :class="
                          r.type === 'entry' ? 'badge-green' : 'badge-red'
                        "
                        class="badge whitespace-nowrap"
                      >
                        {{ r.type === "entry" ? "Entrada" : "Salida" }}
                      </span>
                    </td>
                    <td>
                      <p
                        class="text-xs whitespace-nowrap"
                        style="color: var(--text-muted)"
                      >
                        {{ formatDate(r.timestamp) }}
                      </p>
                    </td>
                    <!-- Horario programado -->
                    <td class="hidden lg:table-cell">
                      <template
                        v-if="
                          (r.type === 'entry' && r.entry_schedule) ||
                          (r.type === 'exit' && r.exit_schedule)
                        "
                      >
                        <span
                          class="text-xs whitespace-nowrap"
                          style="color: var(--text-muted)"
                        >
                          {{
                            r.type === "entry"
                              ? r.entry_schedule
                              : r.exit_schedule
                          }}
                        </span>
                      </template>
                      <span
                        v-else
                        class="text-xs"
                        style="color: var(--text-dim)"
                        >Sin horario</span
                      >
                    </td>
                    <!-- Retraso / Adelanto (calculado en frontend con zona horaria local) -->
                    <td class="hidden lg:table-cell">
                      <template v-if="computeDeviation(r) !== null">
                        <span
                          v-if="computeDeviation(r) === 0"
                          class="badge badge-blue text-xs"
                          >Puntual</span
                        >
                        <span
                          v-else-if="r.type === 'entry'"
                          :class="
                            computeDeviation(r) > 0
                              ? 'badge-red'
                              : 'badge-green'
                          "
                          class="badge text-xs whitespace-nowrap"
                        >
                          {{ formatDeviation(computeDeviation(r), "entry") }}
                        </span>
                        <span
                          v-else
                          :class="
                            computeDeviation(r) < 0
                              ? 'badge-red'
                              : 'badge-green'
                          "
                          class="badge text-xs whitespace-nowrap"
                        >
                          {{ formatDeviation(computeDeviation(r), "exit") }}
                        </span>
                      </template>
                      <span
                        v-else
                        class="text-xs"
                        style="color: var(--text-dim)"
                        >—</span
                      >
                    </td>
                    <td class="hidden md:table-cell">
                      <div class="flex flex-col gap-1">
                        <span
                          :class="
                            r.location_count > 0 ? 'badge-blue' : 'badge-gray'
                          "
                          class="badge"
                        >
                          {{ r.location_count }} pts
                        </span>
                        <span
                          v-if="r.gps_address"
                          class="text-xs truncate max-w-[180px] flex items-center gap-1"
                          style="color: var(--text-muted)"
                          :title="r.gps_address"
                        >
                          <MapPinIcon
                            class="w-3 h-3 flex-shrink-0 text-emerald-500"
                          />
                          {{ r.gps_address }}
                        </span>
                      </div>
                    </td>
                    <!-- Comentario del día -->
                    <td class="hidden xl:table-cell max-w-[200px]">
                      <div class="flex items-center gap-1.5">
                        <span
                          v-if="r.day_comment"
                          class="text-xs truncate"
                          style="color: var(--text-muted)"
                          >{{ r.day_comment }}</span
                        >
                        <span
                          v-else
                          class="text-xs"
                          style="color: var(--text-dim)"
                          >—</span
                        >
                        <button
                          @click="openCommentEditor(r)"
                          class="flex-shrink-0 p-0.5 rounded hover:bg-white/10 transition-colors"
                          title="Editar comentario"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-3.5 h-3.5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            style="color: var(--text-dim)"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"
                            />
                          </svg>
                        </button>
                      </div>
                    </td>
                    <td>
                      <div class="flex items-center gap-2">
                        <button
                          v-if="r.has_site_photo || r.has_selfie_photo"
                          @click="openPhotosModal(r)"
                          class="text-xs px-2 py-1 rounded-lg font-semibold transition-all flex-shrink-0 bg-violet-100 text-violet-700 hover:bg-violet-200 dark:bg-violet-500/15 dark:text-violet-400"
                        >
                          Fotos
                        </button>
                        <span
                          v-else
                          class="text-xs"
                          style="color: var(--text-dim)"
                          >Sin fotos</span
                        >
                        <button
                          @click="openRouteModal(r)"
                          :disabled="r.location_count === 0"
                          class="text-xs px-2 py-1 rounded-lg font-semibold transition-all flex-shrink-0 bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-500/15 dark:text-indigo-400 disabled:opacity-40 disabled:cursor-not-allowed"
                        >
                          Mapa
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Records Pagination -->
            <div
              v-if="totalRecordsPages > 1"
              class="px-5 py-3 flex items-center justify-between"
              style="
                border-top: 1px solid var(--border-subtle);
                background: var(--surface);
              "
            >
              <p class="text-xs text-slate-500">
                Página {{ recordsPage }} de {{ totalRecordsPages }}
              </p>
              <div class="flex items-center gap-2">
                <button
                  class="btn-secondary text-xs px-3 py-1.5"
                  :disabled="recordsPage === 1"
                  @click="recordsPage--"
                >
                  Anterior
                </button>
                <button
                  class="btn-secondary text-xs px-3 py-1.5"
                  :disabled="recordsPage === totalRecordsPages"
                  @click="recordsPage++"
                >
                  Siguiente
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- ===== OPEN SESSIONS TAB ===== -->
        <div v-if="activeTab === 'sessions'" class="space-y-5 animate-in">
          <div class="glass-card overflow-hidden">
            <div
              class="px-5 py-4 flex items-center justify-between"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <div>
                <h3 class="font-semibold" style="color: var(--text)">
                  Sesiones sin cerrar
                </h3>
                <p class="text-xs mt-0.5" style="color: var(--text-muted)">
                  Entradas de las últimas 24 h sin salida correspondiente
                </p>
              </div>
              <button
                @click="loadOpenSessions"
                class="btn-secondary btn-sm flex items-center gap-1.5"
                :disabled="loadingOpenSessions"
              >
                <span
                  v-if="loadingOpenSessions"
                  class="w-3 h-3 border border-brand-500/40 border-t-brand-500 rounded-full animate-spin"
                ></span>
                Actualizar
              </button>
            </div>

            <div
              v-if="openSessions.length === 0 && !loadingOpenSessions"
              class="text-center py-14"
            >
              <CheckCircleIcon
                class="w-12 h-12 mx-auto mb-3 opacity-60"
                style="color: #10b981"
              />
              <p class="font-medium" style="color: #10b981">Todo en orden</p>
              <p class="text-sm mt-1" style="color: var(--text-dim)">
                No hay sesiones abiertas en este momento
              </p>
            </div>

            <div v-else class="overflow-x-auto custom-scroll">
              <table class="data-table w-full">
                <thead>
                  <tr>
                    <th>Usuario</th>
                    <th class="hidden sm:table-cell">Proyecto</th>
                    <th>Hora de entrada</th>
                    <th>Duración</th>
                    <th>Estado</th>
                    <th>Acción</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="s in openSessions" :key="s.record_id">
                    <td>
                      <p class="font-medium text-sm" style="color: var(--text)">
                        {{ toTitleCase(s.first_name) }}
                        {{ toTitleCase(s.last_name) }}
                      </p>
                      <p class="text-xs" style="color: var(--text-muted)">
                        {{ s.email }}
                      </p>
                    </td>
                    <td class="hidden sm:table-cell">
                      <span class="text-xs" style="color: var(--text-muted)">{{
                        s.project_name
                      }}</span>
                    </td>
                    <td>
                      <p
                        class="text-xs whitespace-nowrap"
                        style="color: var(--text-muted)"
                      >
                        {{ formatDate(s.entry_time) }}
                      </p>
                    </td>
                    <td>
                      <span class="badge badge-gray text-xs whitespace-nowrap">
                        {{ formatDuration(s.duration_min) }}
                      </span>
                    </td>
                    <td>
                      <span
                        v-if="s.auto_expired"
                        class="badge badge-gray text-xs"
                        >Expirada (auto)</span
                      >
                      <span
                        v-else
                        class="badge text-xs"
                        style="
                          background: rgba(245, 158, 11, 0.12);
                          border-color: rgba(245, 158, 11, 0.35);
                          color: #d97706;
                        "
                        >Activa — bloquea al usuario</span
                      >
                    </td>
                    <td>
                      <button
                        @click="closeConfirm = { show: true, session: s }"
                        class="text-xs px-2 py-1 rounded-lg font-semibold transition-all bg-rose-100 text-rose-700 hover:bg-rose-200 dark:bg-rose-500/15 dark:text-rose-400"
                      >
                        Cerrar sesión
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Info note -->
          <div
            class="glass-card p-4 flex items-start gap-3"
            style="border-color: rgba(99, 102, 241, 0.2)"
          >
            <InformationCircleIcon
              class="w-5 h-5 text-brand-400 flex-shrink-0 mt-0.5"
            />
            <p class="text-sm" style="color: var(--text-muted)">
              <span class="font-semibold" style="color: var(--text)"
                >¿Cuándo aparece una sesión aquí?</span
              >
              Cuando un usuario registró su entrada pero olvidó registrar la
              salida. Las sesiones con más de
              <strong>20 horas</strong> se marcan como
              <span style="color: var(--text-muted); font-weight: 600"
                >Expiradas</span
              >
              y el sistema ya les permite volver a registrar entrada
              automáticamente. Las sesiones
              <span style="color: #d97706; font-weight: 600">Activas</span>
              aún bloquean al usuario — usa el botón &ldquo;Cerrar sesión&rdquo;
              para desbloquearlo.
            </p>
          </div>
        </div>

        <!-- ===== USERS TAB ===== -->
        <div v-if="activeTab === 'users'" class="space-y-5 animate-in">
          <!-- Filter + actions bar -->
          <div class="glass-card p-4">
            <div class="flex flex-wrap gap-3 items-center">
              <select
                v-model="filters.project"
                @change="loadUsers"
                class="input text-sm py-2 w-52"
              >
                <option value="">Todos los proyectos</option>
                <option v-for="p in projects" :key="p" :value="p">
                  {{ p }}
                </option>
              </select>
              <div class="flex items-center gap-2 flex-1 min-w-[220px]">
                <UserCircleIcon
                  class="w-4 h-4 flex-shrink-0"
                  style="color: var(--text-muted)"
                />
                <input
                  v-model="usersNameFilter"
                  @input="debouncedLoadUsers"
                  type="text"
                  class="input text-sm py-2 flex-1"
                  placeholder="Filtrar por nombre..."
                />
              </div>
              <button
                @click="openProjectsModal"
                class="btn-secondary btn-sm flex items-center gap-1.5"
              >
                <FolderPlusIcon class="w-4 h-4" />
                Proyectos
              </button>
              <p class="text-slate-400 text-sm ml-auto">
                <span style="color: var(--text); font-weight: 600">{{
                  users.length
                }}</span>
                <span style="color: var(--text-muted)"
                  >&nbsp;usuarios registrados</span
                >
              </p>
            </div>
          </div>

          <!-- Users grid -->
          <div v-if="users.length === 0" class="glass-card p-12 text-center">
            <UserCircleIcon
              class="w-12 h-12 mx-auto mb-3 opacity-30"
              style="color: var(--text-muted)"
            />
            <p style="color: var(--text-muted)">Sin usuarios registrados</p>
          </div>
          <div
            v-else
            class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4"
          >
            <div
              v-for="(user, i) in paginatedUsers"
              :key="user.id"
              class="glass-card p-5 animate-in"
              :style="`animation-delay: ${Math.min(i * 0.03, 0.3)}s`"
            >
              <div class="flex items-start gap-3 mb-4">
                <div
                  class="w-12 h-12 rounded-2xl overflow-hidden flex items-center justify-center text-lg font-bold flex-shrink-0"
                  style="
                    background: linear-gradient(
                      135deg,
                      rgba(99, 102, 241, 0.3),
                      rgba(139, 92, 246, 0.3)
                    );
                    border: 1px solid rgba(99, 102, 241, 0.25);
                  "
                >
                  <img
                    v-if="user.avatar_url"
                    :src="user.avatar_url"
                    loading="lazy"
                    class="w-full h-full object-cover"
                  />
                  <span v-else>{{ userInitials(user) }}</span>
                </div>
                <div class="min-w-0">
                  <p class="font-semibold truncate" style="color: var(--text)">
                    {{ toTitleCase(user.first_name) }}
                    {{ toTitleCase(user.last_name) }}
                  </p>
                  <p class="text-xs truncate" style="color: var(--text-muted)">
                    {{ user.email }}
                  </p>
                </div>
              </div>
              <div class="space-y-2">
                <div class="flex items-center gap-2">
                  <span class="text-xs w-16" style="color: var(--text-dim)"
                    >Proyecto</span
                  >
                  <span
                    class="text-xs truncate uppercase"
                    style="color: var(--text)"
                    >{{ user.project_name }}</span
                  >
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-16" style="color: var(--text-dim)"
                    >Registros</span
                  >
                  <span class="badge-blue">{{ user.total_checks }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-16" style="color: var(--text-dim)"
                    >Registro</span
                  >
                  <span class="text-xs" style="color: var(--text-muted)">{{
                    formatDateShort(user.created_at)
                  }}</span>
                </div>
              </div>
              <div
                class="mt-4 pt-3"
                style="border-top: 1px solid var(--border-subtle)"
              >
                <div class="flex gap-2 flex-wrap">
                  <button
                    @click="openEditUserModal(user)"
                    class="flex-1 btn-secondary btn-sm flex items-center justify-center gap-1.5"
                  >
                    <PencilSquareIcon class="w-4 h-4" />
                    Editar
                  </button>
                  <button
                    @click="openScheduleModal(user)"
                    class="flex-1 btn-secondary btn-sm flex items-center justify-center gap-1.5"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="w-4 h-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    Horario
                  </button>
                  <button
                    @click="openPwdModal(user)"
                    class="flex-1 btn-secondary btn-sm flex items-center justify-center gap-1.5"
                  >
                    <LockClosedIcon class="w-4 h-4" />
                    Contraseña
                  </button>
                  <button
                    @click="deleteUserConfirm = { show: true, user }"
                    class="flex-1 btn-danger btn-sm flex items-center justify-center gap-1.5"
                  >
                    <TrashIcon class="w-4 h-4" />
                    Eliminar
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Users Pagination -->
          <div
            v-if="totalUsersPages > 1"
            class="flex items-center justify-between mt-4"
          >
            <p class="text-xs text-slate-500">
              Página {{ usersPage }} de {{ totalUsersPages }}
            </p>
            <div class="flex items-center gap-2">
              <button
                class="btn-secondary text-xs px-3 py-1.5"
                :disabled="usersPage === 1"
                @click="usersPage--"
              >
                Anterior
              </button>
              <button
                class="btn-secondary text-xs px-3 py-1.5"
                :disabled="usersPage === totalUsersPages"
                @click="usersPage++"
              >
                Siguiente
              </button>
            </div>
          </div>
        </div>

        <!-- ===== MISSING RECORDS TAB ===== -->
        <div v-if="activeTab === 'missing'" class="space-y-5 animate-in">
          <!-- Filters -->
          <div class="glass-card p-4">
            <div class="flex flex-wrap gap-3 items-center">
              <input
                type="date"
                v-model="missingFilters.date_from"
                @change="loadMissingRecords"
                class="input text-sm py-2 w-40"
              />
              <span style="color: var(--text-muted)" class="text-sm">a</span>
              <input
                type="date"
                v-model="missingFilters.date_to"
                @change="loadMissingRecords"
                class="input text-sm py-2 w-40"
              />
              <select
                v-model="missingFilters.missing_type"
                @change="loadMissingRecords"
                class="input text-sm py-2 w-44"
              >
                <option value="both">Sin entrada y sin salida</option>
                <option value="entry">Sin entrada</option>
                <option value="exit">Sin salida (tienen entrada)</option>
              </select>
              <button @click="setMissingToday" class="btn-secondary btn-sm">
                Hoy
              </button>
              <button @click="setMissingWeek" class="btn-secondary btn-sm">
                Esta semana
              </button>
              <button
                @click="downloadMissingExcel"
                class="btn-secondary btn-sm flex items-center gap-1.5 ml-auto"
                title="Descargar Excel"
              >
                <ArrowDownTrayIcon class="w-4 h-4" />
                Excel
              </button>
            </div>
          </div>

          <!-- Spinner -->
          <div v-if="loadingMissing" class="glass-card p-10 text-center">
            <div
              class="w-8 h-8 border-4 border-brand-500/30 border-t-brand-500 rounded-full animate-spin mx-auto"
            ></div>
          </div>

          <template v-else>
            <!-- Sin entrada -->
            <div
              v-if="missingFilters.missing_type !== 'exit'"
              class="space-y-2"
            >
              <div class="flex items-center gap-2 px-1">
                <ExclamationCircleIcon class="w-5 h-5 text-rose-500" />
                <h3 class="font-semibold text-sm" style="color: var(--text)">
                  Sin entrada
                </h3>
                <span class="badge-red ml-1">{{
                  missingRecords.missing_entry.length
                }}</span>
              </div>
              <div
                v-if="missingRecords.missing_entry.length === 0"
                class="glass-card p-6 text-center text-sm"
                style="color: var(--text-muted)"
              >
                Todos los usuarios registraron entrada ✓
              </div>
              <div v-else class="glass-card overflow-hidden">
                <table class="data-table w-full">
                  <thead>
                    <tr>
                      <th>Nombre</th>
                      <th class="hidden sm:table-cell">Proyecto</th>
                      <th class="hidden md:table-cell">Correo</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="u in missingRecords.missing_entry"
                      :key="u.user_id"
                    >
                      <td>
                        <span class="font-medium" style="color: var(--text)">
                          {{ toTitleCase(u.first_name) }}
                          {{ toTitleCase(u.last_name) }}
                        </span>
                      </td>
                      <td class="hidden sm:table-cell">
                        <span
                          class="text-xs uppercase"
                          style="color: var(--text-muted)"
                          >{{ u.project_name }}</span
                        >
                      </td>
                      <td class="hidden md:table-cell">
                        <span
                          class="text-xs"
                          style="color: var(--text-muted)"
                          >{{ u.email }}</span
                        >
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Sin salida -->
            <div
              v-if="missingFilters.missing_type !== 'entry'"
              class="space-y-2"
            >
              <div class="flex items-center gap-2 px-1">
                <ExclamationTriangleIcon class="w-5 h-5 text-amber-500" />
                <h3 class="font-semibold text-sm" style="color: var(--text)">
                  Sin salida (tienen entrada)
                </h3>
                <span class="badge-amber ml-1">{{
                  missingRecords.missing_exit.length
                }}</span>
              </div>
              <div
                v-if="missingRecords.missing_exit.length === 0"
                class="glass-card p-6 text-center text-sm"
                style="color: var(--text-muted)"
              >
                Todos los usuarios con entrada registraron salida ✓
              </div>
              <div v-else class="glass-card overflow-hidden">
                <table class="data-table w-full">
                  <thead>
                    <tr>
                      <th>Nombre</th>
                      <th class="hidden sm:table-cell">Proyecto</th>
                      <th class="hidden md:table-cell">Correo</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="u in missingRecords.missing_exit"
                      :key="u.user_id"
                    >
                      <td>
                        <span class="font-medium" style="color: var(--text)">
                          {{ toTitleCase(u.first_name) }}
                          {{ toTitleCase(u.last_name) }}
                        </span>
                      </td>
                      <td class="hidden sm:table-cell">
                        <span
                          class="text-xs uppercase"
                          style="color: var(--text-muted)"
                          >{{ u.project_name }}</span
                        >
                      </td>
                      <td class="hidden md:table-cell">
                        <span
                          class="text-xs"
                          style="color: var(--text-muted)"
                          >{{ u.email }}</span
                        >
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </template>
        </div>

        <!-- ===== PROFILE TAB ===== -->
        <div v-if="activeTab === 'profile'" class="animate-in">
          <!-- Layout: avatar card full-width arriba, luego 2 columnas en lg -->
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
            <!-- Col izquierda: Avatar + Cover + Info de cuenta -->
            <div class="space-y-5">
              <!-- Avatar + Cover card -->
              <div class="glass-card overflow-hidden">
                <!-- Cover -->
                <div
                  class="relative h-36 group cursor-pointer"
                  @click="triggerAdminCoverPick"
                >
                  <div class="absolute inset-0 overflow-hidden">
                    <img
                      v-if="adminCoverPreview || auth.user?.cover_url"
                      :src="adminCoverPreview || auth.user?.cover_url"
                      class="w-full h-full object-cover"
                    />
                    <div
                      v-else
                      class="w-full h-full"
                      style="
                        background: linear-gradient(
                          135deg,
                          #6366f1 0%,
                          #8b5cf6 50%,
                          #06b6d4 100%
                        );
                      "
                    >
                      <div
                        class="absolute inset-0 opacity-20"
                        style="
                          background-image: radial-gradient(
                            #fff 1px,
                            transparent 1px
                          );
                          background-size: 18px 18px;
                        "
                      ></div>
                    </div>
                  </div>
                  <div
                    class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-all duration-200 flex items-center justify-center"
                  >
                    <div
                      class="opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-2 bg-black/60 text-white text-xs font-semibold px-3 py-1.5 rounded-full"
                    >
                      <CameraIcon class="w-3.5 h-3.5" /> Cambiar portada
                    </div>
                  </div>
                  <input
                    ref="adminCoverInput"
                    type="file"
                    accept="image/*"
                    class="hidden"
                    @change="onAdminCoverChange"
                  />
                </div>

                <div class="px-5 pb-5">
                  <div class="flex items-end gap-4 -mt-9 mb-4">
                    <!-- Avatar -->
                    <div
                      class="relative flex-shrink-0 group w-[80px] h-[80px] cursor-pointer"
                      @click="triggerAdminAvatarPick"
                      style="filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.3))"
                    >
                      <div
                        class="w-[80px] h-[80px] rounded-2xl overflow-hidden ring-4 ring-[var(--card-bg)]"
                        style="
                          background: linear-gradient(135deg, #6366f1, #8b5cf6);
                        "
                      >
                        <img
                          v-if="adminAvatarPreview || auth.user?.avatar_url"
                          :src="adminAvatarPreview || auth.user?.avatar_url"
                          class="w-full h-full object-cover"
                        />
                        <div
                          v-else
                          class="w-full h-full flex items-center justify-center text-white text-2xl font-black select-none"
                        >
                          {{ adminInitials }}
                        </div>
                      </div>
                      <div
                        class="absolute inset-0 rounded-2xl bg-black/0 group-hover:bg-black/50 transition-all duration-200 flex items-center justify-center"
                      >
                        <CameraIcon
                          class="w-5 h-5 text-white opacity-0 group-hover:opacity-100 transition-opacity"
                        />
                      </div>
                      <input
                        ref="adminFileInput"
                        type="file"
                        accept="image/*"
                        class="hidden"
                        @change="onAdminAvatarChange"
                      />
                    </div>

                    <div class="mb-1 min-w-0 flex-1 pt-10">
                      <p
                        class="font-black text-xl truncate leading-tight"
                        style="color: var(--text)"
                      >
                        {{ auth.user?.first_name }} {{ auth.user?.last_name }}
                      </p>
                      <p
                        class="text-sm truncate"
                        style="color: var(--text-muted)"
                      >
                        {{ auth.user?.email }}
                      </p>
                      <span class="badge badge-blue mt-1.5">Administrador</span>
                    </div>
                  </div>

                  <!-- Cover pending -->
                  <Transition name="fade">
                    <div
                      v-if="adminCoverPreview"
                      class="flex gap-2 mb-3 p-3 rounded-xl"
                      style="background: var(--input-bg)"
                    >
                      <div
                        class="flex items-center gap-2 flex-1 text-sm"
                        style="color: var(--text-muted)"
                      >
                        <PhotoIcon class="w-4 h-4 flex-shrink-0" /> Nueva
                        portada lista
                      </div>
                      <button
                        @click="saveAdminCover"
                        :disabled="savingAdminCover"
                        class="btn btn-primary btn-sm"
                      >
                        <span
                          v-if="savingAdminCover"
                          class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"
                        ></span>
                        <CheckIcon v-else class="w-3.5 h-3.5" /> Guardar
                      </button>
                      <button
                        @click="cancelAdminCover"
                        class="btn btn-secondary btn-sm px-3"
                      >
                        <XMarkIcon class="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </Transition>

                  <!-- Avatar pending -->
                  <Transition name="fade">
                    <div
                      v-if="adminAvatarPreview"
                      class="flex gap-2 mb-3 p-3 rounded-xl"
                      style="background: var(--input-bg)"
                    >
                      <div
                        class="flex items-center gap-2 flex-1 text-sm"
                        style="color: var(--text-muted)"
                      >
                        <UserCircleIcon class="w-4 h-4 flex-shrink-0" /> Nueva
                        foto de perfil lista
                      </div>
                      <button
                        @click="saveAdminAvatar"
                        :disabled="savingAdminAvatar"
                        class="btn btn-primary btn-sm"
                      >
                        <span
                          v-if="savingAdminAvatar"
                          class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"
                        ></span>
                        <CheckIcon v-else class="w-3.5 h-3.5" /> Guardar
                      </button>
                      <button
                        @click="cancelAdminAvatar"
                        class="btn btn-secondary btn-sm px-3"
                      >
                        <XMarkIcon class="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </Transition>

                  <Transition name="fade">
                    <div
                      v-if="adminAvatarOk"
                      class="flex items-center gap-2 bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 rounded-xl px-4 py-2.5 text-sm font-semibold"
                    >
                      <CheckCircleIcon class="w-4 h-4 flex-shrink-0" /> Foto
                      actualizada correctamente
                    </div>
                  </Transition>
                  <Transition name="fade">
                    <div
                      v-if="adminAvatarErr"
                      class="flex items-center gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-500 rounded-xl px-4 py-2.5 text-sm font-semibold"
                    >
                      <ExclamationTriangleIcon class="w-4 h-4 flex-shrink-0" />
                      {{ adminAvatarErr }}
                    </div>
                  </Transition>
                </div>
              </div>

              <!-- Account info -->
              <div class="glass-card p-6">
                <div class="flex items-center gap-3 mb-5">
                  <div
                    class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                    style="
                      background: rgba(6, 182, 212, 0.12);
                      border: 1px solid rgba(6, 182, 212, 0.2);
                    "
                  >
                    <InformationCircleIcon class="w-5 h-5 text-cyan-400" />
                  </div>
                  <h2 class="font-bold text-base" style="color: var(--text)">
                    Información de cuenta
                  </h2>
                </div>
                <div class="space-y-3">
                  <div
                    class="flex items-center justify-between py-2.5 px-3 rounded-xl"
                    style="background: var(--input-bg)"
                  >
                    <span class="text-sm" style="color: var(--text-muted)"
                      >Nombre</span
                    >
                    <span
                      class="text-sm font-semibold"
                      style="color: var(--text)"
                      >{{ auth.user?.first_name }}
                      {{ auth.user?.last_name }}</span
                    >
                  </div>
                  <div
                    class="flex items-center justify-between py-2.5 px-3 rounded-xl"
                    style="background: var(--input-bg)"
                  >
                    <span class="text-sm" style="color: var(--text-muted)"
                      >Correo</span
                    >
                    <span
                      class="text-sm font-semibold truncate max-w-[200px]"
                      style="color: var(--text)"
                      >{{ auth.user?.email }}</span
                    >
                  </div>
                  <div
                    class="flex items-center justify-between py-2.5 px-3 rounded-xl"
                    style="background: var(--input-bg)"
                  >
                    <span class="text-sm" style="color: var(--text-muted)"
                      >Miembro desde</span
                    >
                    <span
                      class="text-sm font-semibold"
                      style="color: var(--text)"
                      >{{ adminMemberSince }}</span
                    >
                  </div>
                  <div
                    class="flex items-center justify-between py-2.5 px-3 rounded-xl"
                    style="background: var(--input-bg)"
                  >
                    <span class="text-sm" style="color: var(--text-muted)"
                      >Rol</span
                    >
                    <span class="badge badge-blue">Administrador</span>
                  </div>
                </div>
              </div>

              <!-- Appearance card -->
              <div class="glass-card p-6">
                <div class="flex items-center gap-3 mb-5">
                  <div
                    class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                    style="
                      background: rgba(245, 158, 11, 0.12);
                      border: 1px solid rgba(245, 158, 11, 0.2);
                    "
                  >
                    <SunIcon class="w-5 h-5 text-amber-400" />
                  </div>
                  <h2 class="font-bold text-base" style="color: var(--text)">
                    Apariencia
                  </h2>
                </div>
                <div
                  class="flex items-center justify-between py-3 px-4 rounded-xl"
                  style="background: var(--input-bg)"
                >
                  <div>
                    <p class="text-sm font-semibold" style="color: var(--text)">
                      Modo Oscuro
                    </p>
                    <p class="text-xs mt-0.5" style="color: var(--text-muted)">
                      Cambiar tema visual
                    </p>
                  </div>
                  <button
                    @click="theme.toggle()"
                    class="relative inline-flex h-7 w-12 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none"
                    :class="
                      theme.isDark
                        ? 'bg-brand-500'
                        : 'bg-slate-300 dark:bg-surface-600'
                    "
                    role="switch"
                    :aria-checked="theme.isDark"
                  >
                    <span
                      class="pointer-events-none inline-block h-6 w-6 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
                      :class="theme.isDark ? 'translate-x-5' : 'translate-x-0'"
                    ></span>
                  </button>
                </div>
              </div>
            </div>
            <!-- /col izquierda -->

            <!-- Col derecha: Cambiar correo + Cambiar contraseña -->
            <div class="space-y-5">
              <!-- Change email -->
              <div class="glass-card p-6">
                <div class="flex items-center gap-3 mb-5">
                  <div
                    class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                    style="
                      background: rgba(99, 102, 241, 0.12);
                      border: 1px solid rgba(99, 102, 241, 0.2);
                    "
                  >
                    <EnvelopeIcon class="w-5 h-5 text-brand-400" />
                  </div>
                  <div>
                    <h2 class="font-bold text-base" style="color: var(--text)">
                      Cambiar correo
                    </h2>
                    <p class="text-xs" style="color: var(--text-muted)">
                      Actual: {{ auth.user?.email }}
                    </p>
                  </div>
                </div>
                <form
                  @submit.prevent="handleAdminEmailUpdate"
                  class="space-y-4"
                >
                  <div>
                    <label class="input-label">Nuevo correo electrónico</label>
                    <input
                      v-model="adminEmailForm.newEmail"
                      type="email"
                      class="input"
                      placeholder="nuevo@correo.com"
                      required
                      autocomplete="email"
                    />
                  </div>
                  <Transition name="fade">
                    <div
                      v-if="adminEmailMsg.text"
                      class="flex items-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold"
                      :class="
                        adminEmailMsg.ok
                          ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-500'
                          : 'bg-rose-500/10 border border-rose-500/20 text-rose-500'
                      "
                    >
                      <CheckCircleIcon
                        v-if="adminEmailMsg.ok"
                        class="w-4 h-4 flex-shrink-0"
                      />
                      <ExclamationTriangleIcon
                        v-else
                        class="w-4 h-4 flex-shrink-0"
                      />
                      {{ adminEmailMsg.text }}
                    </div>
                  </Transition>
                  <button
                    type="submit"
                    :disabled="savingAdminEmail"
                    class="btn btn-primary btn-md w-full"
                  >
                    <span
                      v-if="savingAdminEmail"
                      class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
                    ></span>
                    <CheckIcon v-else class="w-4 h-4" />
                    {{
                      savingAdminEmail ? "Guardando..." : "Actualizar correo"
                    }}
                  </button>
                </form>
              </div>

              <!-- Change password -->
              <div class="glass-card p-6">
                <div class="flex items-center gap-3 mb-5">
                  <div
                    class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                    style="
                      background: rgba(139, 92, 246, 0.12);
                      border: 1px solid rgba(139, 92, 246, 0.2);
                    "
                  >
                    <LockClosedIcon class="w-5 h-5 text-violet-400" />
                  </div>
                  <div>
                    <h2 class="font-bold text-base" style="color: var(--text)">
                      Cambiar contraseña
                    </h2>
                    <p class="text-xs" style="color: var(--text-muted)">
                      Elige una contraseña segura
                    </p>
                  </div>
                </div>
                <form @submit.prevent="handleAdminPwdUpdate" class="space-y-4">
                  <div>
                    <label class="input-label">Contraseña actual</label>
                    <div class="relative">
                      <input
                        v-model="adminPwdForm.current"
                        :type="showAdminCurrent ? 'text' : 'password'"
                        class="input pr-12"
                        placeholder="••••••••"
                        required
                      />
                      <button
                        type="button"
                        @click="showAdminCurrent = !showAdminCurrent"
                        class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg transition-colors"
                        style="color: var(--text-muted)"
                      >
                        <EyeSlashIcon v-if="showAdminCurrent" class="w-4 h-4" />
                        <EyeIcon v-else class="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                  <div>
                    <label class="input-label">Nueva contraseña</label>
                    <div class="relative">
                      <input
                        v-model="adminPwdForm.newPwd"
                        :type="showAdminNew ? 'text' : 'password'"
                        class="input pr-12"
                        placeholder="Mínimo 8 caracteres"
                        required
                        minlength="8"
                      />
                      <button
                        type="button"
                        @click="showAdminNew = !showAdminNew"
                        class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg transition-colors"
                        style="color: var(--text-muted)"
                      >
                        <EyeSlashIcon v-if="showAdminNew" class="w-4 h-4" />
                        <EyeIcon v-else class="w-4 h-4" />
                      </button>
                    </div>
                    <div
                      v-if="adminPwdForm.newPwd"
                      class="mt-2 flex items-center gap-2"
                    >
                      <div class="flex gap-1 flex-1">
                        <div
                          v-for="n in 4"
                          :key="n"
                          class="h-1 flex-1 rounded-full transition-all duration-300"
                          :style="
                            n <= adminPwdStrength.level
                              ? `background: ${adminPwdStrength.color}`
                              : 'background: var(--input-border)'
                          "
                        ></div>
                      </div>
                      <span
                        class="text-xs font-semibold"
                        :style="`color: ${adminPwdStrength.color}`"
                        >{{ adminPwdStrength.label }}</span
                      >
                    </div>
                  </div>
                  <div>
                    <label class="input-label"
                      >Confirmar nueva contraseña</label
                    >
                    <input
                      v-model="adminPwdForm.confirm"
                      type="password"
                      class="input"
                      :class="
                        adminPwdForm.confirm && adminPwdMismatch
                          ? 'ring-2 ring-rose-500/30 border-rose-500/40'
                          : ''
                      "
                      placeholder="••••••••"
                      required
                    />
                    <p
                      v-if="adminPwdForm.confirm && adminPwdMismatch"
                      class="text-xs text-rose-400 mt-1 font-semibold flex items-center gap-1"
                    >
                      <ExclamationTriangleIcon class="w-3.5 h-3.5" /> Las
                      contraseñas no coinciden
                    </p>
                  </div>
                  <Transition name="fade">
                    <div
                      v-if="adminPwdMsg.text"
                      class="flex items-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold"
                      :class="
                        adminPwdMsg.ok
                          ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-500'
                          : 'bg-rose-500/10 border border-rose-500/20 text-rose-500'
                      "
                    >
                      <CheckCircleIcon
                        v-if="adminPwdMsg.ok"
                        class="w-4 h-4 flex-shrink-0"
                      />
                      <ExclamationTriangleIcon
                        v-else
                        class="w-4 h-4 flex-shrink-0"
                      />
                      {{ adminPwdMsg.text }}
                    </div>
                  </Transition>
                  <button
                    type="submit"
                    :disabled="savingAdminPwd || adminPwdMismatch"
                    class="btn w-full btn-md text-white font-bold"
                    style="
                      background: linear-gradient(135deg, #7c3aed, #6366f1);
                      box-shadow: 0 4px 15px rgba(124, 58, 237, 0.25);
                    "
                  >
                    <span
                      v-if="savingAdminPwd"
                      class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
                    ></span>
                    <LockClosedIcon v-else class="w-4 h-4" />
                    {{ savingAdminPwd ? "Guardando..." : "Cambiar contraseña" }}
                  </button>
                </form>
              </div>
            </div>
            <!-- /col derecha -->
          </div>
          <!-- /grid -->
        </div>
      </div>
    </main>

    <!-- Mobile bottom nav -->
    <nav
      class="lg:hidden fixed bottom-0 left-0 right-0 z-20 px-2 py-2 safe-bottom flex justify-around"
      style="
        background: var(--nav-bg);
        backdrop-filter: blur(12px);
        border-top: 1px solid var(--border-subtle);
      "
    >
      <button
        v-for="item in navItems"
        :key="item.id"
        @click="activeTab = item.id"
        :class="
          activeTab === item.id
            ? 'text-brand-500 dark:text-brand-400'
            : 'text-slate-500 dark:text-slate-500'
        "
        class="flex flex-col items-center gap-1 px-4 py-1 transition-colors"
      >
        <component :is="item.icon" class="w-6 h-6" />
        <span class="text-xs font-semibold">{{ item.label }}</span>
      </button>
    </nav>

    <!-- ===== MODALS ===== -->
    <Teleport to="body">
      <!-- Route modal -->
      <Transition name="modal">
        <div
          v-if="routeModal.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-3 py-4"
          style="background: rgba(0, 0, 0, 0.88); backdrop-filter: blur(10px)"
        >
          <div
            class="w-full max-w-3xl glass-card flex flex-col animate-in overflow-hidden"
            style="
              background: var(--modal-bg);
              max-height: 94vh;
              border: 1px solid rgba(99, 102, 241, 0.2);
            "
          >
            <!-- Header premium -->
            <div
              class="flex items-center justify-between px-5 py-4 flex-shrink-0"
              style="
                background: linear-gradient(
                  135deg,
                  rgba(99, 102, 241, 0.15),
                  rgba(139, 92, 246, 0.1)
                );
                border-bottom: 1px solid var(--border-subtle);
              "
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                  style="background: linear-gradient(135deg, #6366f1, #8b5cf6)"
                >
                  <MapPinIcon class="w-5 h-5 text-white" />
                </div>
                <div>
                  <h3 class="font-bold text-base" style="color: var(--text)">
                    Recorrido GPS
                  </h3>
                  <p class="text-xs mt-0.5" style="color: var(--text-muted)">
                    {{ routeModal.user }}
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <a
                  v-if="routeModal.points.length > 0"
                  :href="googleMapsUrl(routeModal.points)"
                  target="_blank"
                  rel="noopener"
                  class="btn-secondary btn-sm flex items-center gap-1.5 text-xs"
                >
                  <MapPinIcon class="w-3.5 h-3.5" /> Ver en Google Maps
                </a>
                <button
                  @click="
                    routeModal.show = false;
                    destroyMap();
                  "
                  class="w-8 h-8 rounded-lg flex items-center justify-center transition-all hover:bg-white/10"
                  style="color: var(--text-muted)"
                >
                  <XMarkIcon class="w-5 h-5" />
                </button>
              </div>
            </div>

            <!-- Stats bar -->
            <div
              v-if="routeModal.points.length > 0"
              class="flex items-center gap-4 px-5 py-2.5 flex-shrink-0 flex-wrap"
              style="
                background: rgba(99, 102, 241, 0.05);
                border-bottom: 1px solid var(--border-subtle);
              "
            >
              <div class="flex items-center gap-1.5">
                <div
                  class="w-3 h-3 rounded-full bg-emerald-400 flex-shrink-0"
                ></div>
                <span class="text-xs" style="color: var(--text-muted)"
                  >Inicio:</span
                >
                <span
                  class="text-xs font-semibold"
                  style="color: var(--text)"
                  >{{ formatTimeOnly(routeModal.points[0]?.recorded_at) }}</span
                >
              </div>
              <div class="flex items-center gap-1.5">
                <div
                  class="w-3 h-3 rounded-full bg-rose-400 flex-shrink-0"
                ></div>
                <span class="text-xs" style="color: var(--text-muted)"
                  >Fin:</span
                >
                <span
                  class="text-xs font-semibold"
                  style="color: var(--text)"
                  >{{
                    formatTimeOnly(
                      routeModal.points[routeModal.points.length - 1]
                        ?.recorded_at,
                    )
                  }}</span
                >
              </div>
              <div class="flex items-center gap-1.5">
                <SignalIcon class="w-3.5 h-3.5 text-brand-400" />
                <span class="text-xs font-semibold" style="color: var(--text)"
                  >{{ routeModal.points.length }} puntos GPS</span
                >
              </div>
              <div v-if="routeModal.duration" class="flex items-center gap-1.5">
                <span class="text-xs" style="color: var(--text-muted)"
                  >Duración:</span
                >
                <span
                  class="text-xs font-semibold"
                  style="color: var(--text)"
                  >{{ routeModal.duration }}</span
                >
              </div>
            </div>

            <!-- Map or empty -->
            <div
              v-if="routeModal.points.length === 0"
              class="text-center py-20 text-sm"
              style="color: var(--text-muted)"
            >
              <MapPinIcon class="w-10 h-10 mx-auto mb-3 opacity-30" />
              Sin puntos GPS registrados
            </div>
            <div v-else class="flex-1 overflow-auto">
              <div ref="mapContainer" style="height: 420px; width: 100%"></div>

              <!-- Address section -->
              <div
                class="px-5 py-4"
                style="border-top: 1px solid var(--border-subtle)"
              >
                <div class="flex items-start gap-3">
                  <div
                    class="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0 mt-0.5"
                    style="
                      background: rgba(99, 102, 241, 0.15);
                      border: 1px solid rgba(99, 102, 241, 0.25);
                    "
                  >
                    <MapPinIcon class="w-4 h-4 text-brand-400" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <p
                      class="text-xs font-semibold mb-1"
                      style="color: var(--text-muted)"
                    >
                      📍 Dirección detectada (punto de inicio)
                    </p>
                    <div
                      v-if="routeModal.loadingAddress"
                      class="flex items-center gap-2"
                    >
                      <div
                        class="w-3 h-3 border border-brand-500/40 border-t-brand-500 rounded-full animate-spin"
                      ></div>
                      <span class="text-xs" style="color: var(--text-muted)"
                        >Obteniendo dirección...</span
                      >
                    </div>
                    <p
                      v-else-if="routeModal.address"
                      class="text-sm font-medium"
                      style="color: var(--text)"
                    >
                      {{ routeModal.address }}
                    </p>
                    <p v-else class="text-xs" style="color: var(--text-dim)">
                      Dirección no disponible
                    </p>
                    <p
                      v-if="routeModal.points.length > 0"
                      class="text-xs mt-1"
                      style="color: var(--text-dim)"
                    >
                      Coordenadas:
                      {{ routeModal.points[0].latitude.toFixed(6) }},
                      {{ routeModal.points[0].longitude.toFixed(6) }}
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
        <div
          v-if="suspiciousModal.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
        >
          <div
            class="w-full max-w-md glass-card animate-in"
            style="background: var(--modal-bg)"
          >
            <div
              class="flex items-center justify-between px-5 py-4"
              style="border-bottom: 1px solid rgba(245, 158, 11, 0.3)"
            >
              <div class="flex items-center gap-2">
                <ExclamationTriangleIcon class="w-5 h-5 text-amber-400" />
                <h3 class="font-bold text-amber-400">
                  Alerta de fraude detectado
                </h3>
              </div>
              <button
                @click="suspiciousModal.show = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center"
                style="color: var(--text-muted)"
              >
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-5 space-y-4">
              <div
                class="rounded-xl p-4"
                style="
                  background: rgba(245, 158, 11, 0.08);
                  border: 1px solid rgba(245, 158, 11, 0.25);
                "
              >
                <p
                  class="text-sm font-semibold mb-1"
                  style="color: var(--text)"
                >
                  {{ suspiciousModal.user }}
                </p>
                <p class="text-xs" style="color: var(--text-muted)">
                  {{ suspiciousModal.date }}
                </p>
              </div>
              <div>
                <p class="text-xs font-semibold mb-2 text-amber-400">
                  Razón de la alerta:
                </p>
                <p class="text-sm" style="color: var(--text)">
                  {{ suspiciousModal.reason }}
                </p>
              </div>
              <div
                v-if="suspiciousModal.ipCountry"
                class="flex items-center gap-2"
              >
                <span class="text-xs" style="color: var(--text-dim)"
                  >IP detectada en:</span
                >
                <span class="text-xs font-medium" style="color: var(--text)"
                  >{{
                    suspiciousModal.ipCity ? suspiciousModal.ipCity + ", " : ""
                  }}{{ suspiciousModal.ipCountry }}</span
                >
              </div>
              <div
                class="rounded-xl p-3"
                style="
                  background: rgba(239, 68, 68, 0.07);
                  border: 1px solid rgba(239, 68, 68, 0.2);
                "
              >
                <p class="text-xs" style="color: var(--text-muted)">
                  El registro fue guardado normalmente. Esta alerta es solo
                  informativa para que el administrador investigue.
                </p>
              </div>
              <button
                @click="suspiciousModal.show = false"
                class="btn-secondary btn-md w-full"
              >
                Cerrar
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Photos modal -->
      <Transition name="modal">
        <div
          v-if="photosModal.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.9); backdrop-filter: blur(8px)"
        >
          <div
            class="w-full max-w-lg glass-card animate-in"
            style="background: var(--modal-bg)"
          >
            <div
              class="flex items-center justify-between px-5 py-4"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <div>
                <h3 class="font-bold" style="color: var(--text)">
                  Fotografías
                </h3>
                <p class="text-xs mt-0.5" style="color: var(--text-muted)">
                  {{ photosModal.user }} · {{ photosModal.date }}
                </p>
              </div>
              <button
                @click="photosModal.show = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center transition-all"
                style="color: var(--text-muted)"
              >
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-5 grid grid-cols-2 gap-4">
              <div>
                <p
                  class="text-xs mb-2 flex items-center gap-1"
                  style="color: var(--text-muted)"
                >
                  <BuildingOffice2Icon class="w-3 h-3" /> Fotografía del sitio
                </p>
                <div v-if="photosModal.site" class="rounded-xl overflow-hidden">
                  <img
                    :src="photosModal.site"
                    class="w-full aspect-square object-cover"
                  />
                </div>
                <div
                  v-else
                  class="rounded-xl flex items-center justify-center aspect-square text-slate-600 text-sm"
                  style="
                    background: var(--input-bg);
                    border: 1px dashed var(--input-border);
                  "
                >
                  Sin foto
                </div>
              </div>
              <div>
                <p
                  class="text-xs mb-2 flex items-center gap-1"
                  style="color: var(--text-muted)"
                >
                  <UserCircleIcon class="w-3 h-3" /> Selfie
                </p>
                <div
                  v-if="photosModal.selfie"
                  class="rounded-xl overflow-hidden"
                >
                  <img
                    :src="photosModal.selfie"
                    class="w-full aspect-square object-cover"
                  />
                </div>
                <div
                  v-else
                  class="rounded-xl flex items-center justify-center aspect-square text-slate-600 text-sm"
                  style="
                    background: var(--input-bg);
                    border: 1px dashed var(--input-border);
                  "
                >
                  Sin foto
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Schedule modal -->
      <Transition name="modal">
        <div
          v-if="scheduleModal.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
          @click.self="scheduleModal.show = false"
        >
          <div
            class="w-full max-w-sm glass-card animate-in"
            style="background: var(--modal-bg)"
          >
            <div
              class="flex items-center justify-between px-5 pt-5 pb-3"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <div>
                <h3 class="font-bold" style="color: var(--text)">
                  Horario Laboral
                </h3>
                <p class="text-xs mt-0.5" style="color: var(--text-muted)">
                  {{ scheduleModal.user?.first_name }}
                  {{ scheduleModal.user?.last_name }}
                </p>
              </div>
              <button
                @click="scheduleModal.show = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center"
                style="color: var(--text-muted)"
              >
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-5 space-y-4">
              <!-- Lunes – Viernes -->
              <p
                class="text-xs font-bold uppercase tracking-wider"
                style="color: var(--text-muted)"
              >
                Lunes – Viernes
              </p>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label
                    class="block text-xs font-semibold mb-1.5"
                    style="color: var(--text-muted)"
                    >Entrada</label
                  >
                  <input
                    type="time"
                    v-model="scheduleModal.entry_time"
                    class="input w-full text-sm py-2"
                  />
                </div>
                <div>
                  <label
                    class="block text-xs font-semibold mb-1.5"
                    style="color: var(--text-muted)"
                    >Salida</label
                  >
                  <input
                    type="time"
                    v-model="scheduleModal.exit_time"
                    class="input w-full text-sm py-2"
                  />
                </div>
              </div>
              <!-- Sábado -->
              <div
                style="border-top: 1px solid var(--border-subtle)"
                class="pt-3"
              >
                <p
                  class="text-xs font-bold uppercase tracking-wider mb-3"
                  style="color: var(--text-muted)"
                >
                  Sábado <span class="font-normal normal-case">(opcional)</span>
                </p>
                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label
                      class="block text-xs font-semibold mb-1.5"
                      style="color: var(--text-muted)"
                      >Entrada</label
                    >
                    <input
                      type="time"
                      v-model="scheduleModal.sat_entry_time"
                      class="input w-full text-sm py-2"
                    />
                  </div>
                  <div>
                    <label
                      class="block text-xs font-semibold mb-1.5"
                      style="color: var(--text-muted)"
                      >Salida</label
                    >
                    <input
                      type="time"
                      v-model="scheduleModal.sat_exit_time"
                      class="input w-full text-sm py-2"
                    />
                  </div>
                </div>
              </div>
              <div class="flex gap-3 pt-1">
                <button
                  @click="scheduleModal.show = false"
                  class="btn-secondary flex-1"
                >
                  Cancelar
                </button>
                <button
                  @click="saveSchedule"
                  :disabled="
                    scheduleModal.saving ||
                    !scheduleModal.entry_time ||
                    !scheduleModal.exit_time
                  "
                  class="btn-primary flex-1 disabled:opacity-50"
                >
                  {{ scheduleModal.saving ? "Guardando..." : "Guardar" }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Comment edit modal -->
      <Transition name="modal">
        <div
          v-if="editingComment"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
          @click.self="editingComment = null"
        >
          <div
            class="w-full max-w-md glass-card animate-in"
            style="background: var(--modal-bg)"
          >
            <div
              class="flex items-center justify-between px-5 pt-5 pb-3"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <h3 class="font-bold" style="color: var(--text)">
                Comentario del día · {{ editingComment?.record_date }}
              </h3>
              <button
                @click="editingComment = null"
                class="w-8 h-8 rounded-lg flex items-center justify-center"
                style="color: var(--text-muted)"
              >
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-5 space-y-4">
              <textarea
                v-model="editingComment.text"
                rows="4"
                class="input w-full text-sm py-2 resize-none"
                placeholder="Motivo de llegada tarde, salida anticipada, ausencia..."
              ></textarea>
              <div class="flex gap-3">
                <button
                  @click="editingComment = null"
                  class="btn-secondary flex-1"
                >
                  Cancelar
                </button>
                <button
                  @click="saveComment"
                  :disabled="!editingComment.text.trim()"
                  class="btn-primary flex-1 disabled:opacity-50"
                >
                  Guardar
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Password reset modal -->
      <Transition name="modal">
        <div
          v-if="pwdModal.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
          @click.self="pwdModal.show = false"
        >
          <div
            class="w-full max-w-sm glass-card animate-in"
            style="background: var(--modal-bg)"
          >
            <div
              class="flex items-center justify-between px-5 py-4"
              style="border-bottom: 1px solid rgba(139, 92, 246, 0.3)"
            >
              <div class="flex items-center gap-2">
                <LockClosedIcon class="w-5 h-5 text-violet-400" />
                <div>
                  <h3 class="font-bold text-base" style="color: var(--text)">
                    Cambiar contraseña
                  </h3>
                  <p class="text-xs" style="color: var(--text-muted)">
                    {{ pwdModal.user?.first_name }}
                    {{ pwdModal.user?.last_name }}
                  </p>
                </div>
              </div>
              <button
                @click="pwdModal.show = false"
                class="w-8 h-8 rounded-lg flex items-center justify-center"
                style="color: var(--text-muted)"
              >
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-5 space-y-4">
              <div>
                <label
                  class="block text-xs font-semibold mb-1.5"
                  style="color: var(--text-muted)"
                  >Nueva contraseña</label
                >
                <div class="relative">
                  <input
                    v-model="pwdModal.password"
                    :type="pwdModal.showPwd ? 'text' : 'password'"
                    class="input w-full pr-11"
                    placeholder="Mínimo 8 caracteres"
                    minlength="8"
                    autocomplete="new-password"
                  />
                  <button
                    type="button"
                    @click="pwdModal.showPwd = !pwdModal.showPwd"
                    class="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg transition-colors"
                    style="color: var(--text-muted)"
                  >
                    <EyeSlashIcon v-if="pwdModal.showPwd" class="w-4 h-4" />
                    <EyeIcon v-else class="w-4 h-4" />
                  </button>
                </div>
              </div>
              <div>
                <label
                  class="block text-xs font-semibold mb-1.5"
                  style="color: var(--text-muted)"
                  >Confirmar contraseña</label
                >
                <input
                  v-model="pwdModal.confirm"
                  type="password"
                  class="input w-full"
                  :class="
                    pwdModal.confirm && pwdModal.password !== pwdModal.confirm
                      ? 'ring-2 ring-rose-500/30 border-rose-500/40'
                      : ''
                  "
                  placeholder="Repite la contraseña"
                  autocomplete="new-password"
                />
                <p
                  v-if="
                    pwdModal.confirm && pwdModal.password !== pwdModal.confirm
                  "
                  class="text-xs text-rose-400 mt-1 font-semibold flex items-center gap-1"
                >
                  <ExclamationTriangleIcon class="w-3.5 h-3.5" /> Las
                  contraseñas no coinciden
                </p>
              </div>
              <Transition name="fade">
                <div
                  v-if="pwdModal.msg.text"
                  class="flex items-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold"
                  :class="
                    pwdModal.msg.ok
                      ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-500'
                      : 'bg-rose-500/10 border border-rose-500/20 text-rose-500'
                  "
                >
                  <CheckCircleIcon
                    v-if="pwdModal.msg.ok"
                    class="w-4 h-4 flex-shrink-0"
                  />
                  <ExclamationTriangleIcon
                    v-else
                    class="w-4 h-4 flex-shrink-0"
                  />
                  {{ pwdModal.msg.text }}
                </div>
              </Transition>
              <div class="flex gap-3 pt-1">
                <button
                  @click="pwdModal.show = false"
                  class="btn-secondary flex-1"
                >
                  Cancelar
                </button>
                <button
                  @click="savePwd"
                  :disabled="
                    pwdModal.saving ||
                    pwdModal.password.length < 8 ||
                    pwdModal.password !== pwdModal.confirm
                  "
                  class="flex-1 btn-md text-white font-bold disabled:opacity-50 disabled:cursor-not-allowed rounded-xl flex items-center justify-center gap-2"
                  style="background: linear-gradient(135deg, #7c3aed, #6366f1)"
                >
                  <span
                    v-if="pwdModal.saving"
                    class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
                  ></span>
                  <LockClosedIcon v-else class="w-4 h-4" />
                  {{ pwdModal.saving ? "Guardando..." : "Cambiar contraseña" }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Close session confirmation modal -->
      <Transition name="modal">
        <div
          v-if="closeConfirm.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
          @click.self="closeConfirm = { show: false, session: null }"
        >
          <div
            class="w-full max-w-sm glass-card p-6 animate-in text-center"
            style="background: var(--modal-bg)"
          >
            <div
              class="w-14 h-14 mx-auto mb-4 rounded-full flex items-center justify-center"
              style="background: rgba(239, 68, 68, 0.1)"
            >
              <ClockIcon class="w-7 h-7 text-rose-500" />
            </div>
            <h3 class="font-bold text-lg mb-1" style="color: var(--text)">
              Cerrar sesión manualmente
            </h3>
            <p class="text-sm font-semibold mb-1" style="color: var(--text)">
              {{ closeConfirm.session?.first_name }}
              {{ closeConfirm.session?.last_name }}
            </p>
            <p class="text-sm mb-5" style="color: var(--text-muted)">
              La entrada se marcará como cerrada. El usuario podrá registrar una
              nueva entrada de inmediato.
            </p>
            <div class="flex gap-3">
              <button
                @click="closeConfirm = { show: false, session: null }"
                class="btn-secondary flex-1"
              >
                Cancelar
              </button>
              <button @click="handleCloseSession" class="btn-danger flex-1">
                Cerrar sesión
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Logout confirmation modal -->
      <Transition name="modal">
        <div
          v-if="showLogoutModal"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
        >
          <div
            class="w-full max-w-sm glass-card p-6 animate-in text-center"
            style="background: var(--modal-bg)"
          >
            <div
              class="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center bg-rose-500/10 text-rose-500"
            >
              <ArrowRightOnRectangleIcon class="w-8 h-8" />
            </div>
            <h3 class="font-bold text-lg mb-2" style="color: var(--text)">
              ¿Cerrar sesión?
            </h3>
            <p class="text-sm mb-6" style="color: var(--text-muted)">
              Tendrás que volver a iniciar sesión para acceder al panel
              administrativo.
            </p>
            <div class="flex gap-3">
              <button
                @click="showLogoutModal = false"
                class="btn-secondary flex-1"
              >
                Cancelar
              </button>
              <button @click="handleLogout" class="btn-danger flex-1">
                Salir
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Delete user confirmation modal -->
      <Transition name="modal">
        <div
          v-if="deleteUserConfirm.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
          @click.self="deleteUserConfirm = { show: false, user: null, saving: false }"
        >
          <div
            class="w-full max-w-sm glass-card p-6 animate-in text-center"
            style="background: var(--modal-bg)"
          >
            <div
              class="w-14 h-14 mx-auto mb-4 rounded-full flex items-center justify-center"
              style="background: rgba(239, 68, 68, 0.1)"
            >
              <TrashIcon class="w-7 h-7 text-rose-500" />
            </div>
            <h3 class="font-bold text-lg mb-1" style="color: var(--text)">
              Eliminar usuario
            </h3>
            <p class="text-sm font-semibold mb-1" style="color: var(--text)">
              {{ deleteUserConfirm.user?.first_name }}
              {{ deleteUserConfirm.user?.last_name }}
            </p>
            <p class="text-sm mb-5" style="color: var(--text-muted)">
              Esta acción es permanente y eliminará su historial de registros.
            </p>
            <div class="flex gap-3">
              <button
                @click="deleteUserConfirm = { show: false, user: null, saving: false }"
                class="btn-secondary flex-1"
                :disabled="deleteUserConfirm.saving"
              >
                Cancelar
              </button>
              <button
                @click="handleDeleteUser"
                class="btn-danger flex-1"
                :disabled="deleteUserConfirm.saving"
              >
                {{ deleteUserConfirm.saving ? "Eliminando..." : "Eliminar" }}
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- ── Edit User Modal ──────────────────────────────────────────────── -->
      <Transition name="modal">
        <div
          v-if="editUserModal.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
          @click.self="editUserModal.show = false"
        >
          <div
            class="w-full max-w-md glass-card animate-in"
            style="background: var(--modal-bg)"
          >
            <div
              class="flex items-center justify-between px-5 pt-5 pb-3"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <h3 class="font-bold" style="color: var(--text)">
                Editar usuario
              </h3>
              <button
                @click="editUserModal.show = false"
                class="p-1 rounded-lg hover:bg-white/10"
              >
                <XMarkIcon class="w-5 h-5" style="color: var(--text-muted)" />
              </button>
            </div>
            <div class="p-5 space-y-4">
              <div class="grid grid-cols-2 gap-3">
                <div class="space-y-1">
                  <label class="label-base text-xs">Nombre(s)</label>
                  <input
                    v-model="editUserModal.first_name"
                    type="text"
                    class="input-base text-sm"
                    placeholder="Nombre"
                  />
                </div>
                <div class="space-y-1">
                  <label class="label-base text-xs">Apellido(s)</label>
                  <input
                    v-model="editUserModal.last_name"
                    type="text"
                    class="input-base text-sm"
                    placeholder="Apellido"
                  />
                </div>
              </div>
              <div class="space-y-1">
                <label class="label-base text-xs">Proyecto</label>
                <select
                  v-model="editUserModal.project_name"
                  class="input-base text-sm"
                >
                  <option value="" disabled>Selecciona un proyecto...</option>
                  <option v-for="p in projects" :key="p" :value="p">
                    {{ p }}
                  </option>
                </select>
              </div>
              <div class="space-y-1">
                <label class="label-base text-xs">Correo electrónico</label>
                <input
                  v-model="editUserModal.email"
                  type="email"
                  class="input-base text-sm"
                  placeholder="correo@ejemplo.com"
                />
              </div>
              <Transition name="fade">
                <p
                  v-if="editUserModal.msg.text"
                  :class="
                    editUserModal.msg.ok ? 'text-emerald-500' : 'text-rose-500'
                  "
                  class="text-xs font-medium"
                >
                  {{ editUserModal.msg.text }}
                </p>
              </Transition>
            </div>
            <div class="flex gap-3 px-5 pb-5">
              <button
                @click="editUserModal.show = false"
                class="btn-secondary flex-1"
              >
                Cancelar
              </button>
              <button
                @click="saveEditUser"
                class="btn-primary flex-1"
                :disabled="
                  editUserModal.saving ||
                  !editUserModal.first_name ||
                  !editUserModal.last_name ||
                  !editUserModal.project_name ||
                  !editUserModal.email
                "
              >
                {{ editUserModal.saving ? "Guardando..." : "Guardar" }}
              </button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- ── Projects Management Modal ───────────────────────────────────── -->
      <Transition name="modal">
        <div
          v-if="projectsModal.show"
          class="fixed inset-0 z-50 flex items-center justify-center px-4"
          style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px)"
          @click.self="projectsModal.show = false"
        >
          <div
            class="w-full max-w-md glass-card animate-in"
            style="
              background: var(--modal-bg);
              max-height: 85vh;
              display: flex;
              flex-direction: column;
            "
          >
            <div
              class="flex items-center justify-between px-5 pt-5 pb-3 flex-shrink-0"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <h3 class="font-bold" style="color: var(--text)">
                Gestionar proyectos
              </h3>
              <button
                @click="projectsModal.show = false"
                class="p-1 rounded-lg hover:bg-white/10"
              >
                <XMarkIcon class="w-5 h-5" style="color: var(--text-muted)" />
              </button>
            </div>

            <!-- Add new project -->
            <div
              class="px-5 pt-4 pb-3 flex-shrink-0"
              style="border-bottom: 1px solid var(--border-subtle)"
            >
              <label class="label-base text-xs mb-2 block"
                >Agregar proyecto</label
              >
              <div class="flex gap-2">
                <input
                  v-model="projectsModal.newName"
                  type="text"
                  class="input-base text-sm flex-1 uppercase"
                  placeholder="Nombre del proyecto"
                  @keyup.enter="saveProject"
                  @input="
                    projectsModal.newName = projectsModal.newName.toUpperCase()
                  "
                />
                <button
                  @click="saveProject"
                  class="btn-primary px-4"
                  :disabled="
                    projectsModal.saving || !projectsModal.newName.trim()
                  "
                >
                  <PlusIcon class="w-4 h-4" />
                </button>
              </div>
              <p v-if="projectsModal.error" class="text-xs text-rose-500 mt-1">
                {{ projectsModal.error }}
              </p>
            </div>

            <!-- Projects list -->
            <div class="overflow-y-auto flex-1 p-3">
              <div
                v-if="projects.length === 0"
                class="text-center py-6 text-sm"
                style="color: var(--text-muted)"
              >
                No hay proyectos registrados
              </div>
              <div v-else class="space-y-1">
                <div
                  v-for="p in projects"
                  :key="p"
                  class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-white/5 group"
                >
                  <!-- Editing mode -->
                  <template v-if="projectsModal.editing === p">
                    <input
                      v-model="projectsModal.editName"
                      type="text"
                      class="input-base text-sm flex-1 uppercase py-1"
                      @input="projectsModal.editName = projectsModal.editName.toUpperCase()"
                      @keyup.enter="saveRenameProject(p)"
                      @keyup.escape="projectsModal.editing = null"
                      autofocus
                    />
                    <button
                      @click="saveRenameProject(p)"
                      class="p-1.5 rounded-lg text-emerald-400 hover:bg-emerald-500/10 transition-all"
                      :disabled="projectsModal.saving"
                      title="Guardar"
                    >
                      <CheckIcon class="w-4 h-4" />
                    </button>
                    <button
                      @click="projectsModal.editing = null"
                      class="p-1.5 rounded-lg text-slate-400 hover:bg-white/10 transition-all"
                      title="Cancelar"
                    >
                      <XMarkIcon class="w-4 h-4" />
                    </button>
                  </template>

                  <!-- View mode -->
                  <template v-else>
                    <span
                      class="text-sm font-medium uppercase flex-1"
                      style="color: var(--text)"
                      >{{ p }}</span
                    >
                    <button
                      @click="startEditProject(p)"
                      class="p-1.5 rounded-lg text-brand-400 opacity-0 group-hover:opacity-100 hover:bg-brand-500/10 transition-all"
                      title="Renombrar proyecto"
                    >
                      <PencilSquareIcon class="w-4 h-4" />
                    </button>
                    <button
                      @click="deleteProject(p)"
                      class="p-1.5 rounded-lg text-rose-400 opacity-0 group-hover:opacity-100 hover:bg-rose-500/10 transition-all"
                      title="Eliminar proyecto"
                    >
                      <TrashIcon class="w-4 h-4" />
                    </button>
                  </template>
                </div>
              </div>
            </div>

            <div class="px-5 pb-5 pt-3 flex-shrink-0">
              <button
                @click="projectsModal.show = false"
                class="btn-secondary w-full"
              >
                Cerrar
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import api from "@/api";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import {
  CheckCircleIcon,
  MapPinIcon,
  ArrowRightOnRectangleIcon,
  SignalIcon,
  ClipboardDocumentListIcon,
  UsersIcon,
  ChartBarIcon,
  UserCircleIcon,
  BuildingOffice2Icon,
  CalendarIcon,
  MagnifyingGlassIcon,
  XMarkIcon,
  ExclamationTriangleIcon,
  CameraIcon,
  EnvelopeIcon,
  LockClosedIcon,
  CheckIcon,
  PhotoIcon,
  EyeIcon,
  EyeSlashIcon,
  InformationCircleIcon,
  SunIcon,
  ClockIcon,
  PencilSquareIcon,
  TrashIcon,
  PlusIcon,
  ExclamationCircleIcon,
  FolderPlusIcon,
  ArrowDownTrayIcon,
} from "@heroicons/vue/24/outline";
import { useThemeStore } from "@/stores/theme";
import * as XLSX from "xlsx";

const auth = useAuthStore();
const theme = useThemeStore();
const router = useRouter();

const activeTab = ref("overview");
const navItems = [
  {
    id: "overview",
    icon: ChartBarIcon,
    label: "Resumen",
    description: "Vista general del día",
  },
  {
    id: "records",
    icon: ClipboardDocumentListIcon,
    label: "Registros",
    description: "Todas las entradas y salidas",
  },
  {
    id: "users",
    icon: UsersIcon,
    label: "Usuarios",
    description: "Personas registradas en el sistema",
  },
  {
    id: "sessions",
    icon: ClockIcon,
    label: "Sesiones",
    description: "Entradas sin salida registrada",
  },
  {
    id: "missing",
    icon: ExclamationCircleIcon,
    label: "Sin Registros",
    description: "Usuarios sin entrada o salida en el día",
  },
  {
    id: "profile",
    icon: UserCircleIcon,
    label: "Mi Perfil",
    description: "Configuración de tu cuenta",
  },
];
const currentNavItem = computed(() =>
  navItems.find((n) => n.id === activeTab.value),
);

// Data
const stats = ref({
  total_users: 0,
  records_today: 0,
  entries_total: 0,
  exits_total: 0,
  active_now: 0,
});
const records = ref([]);
const users = ref([]);
const projects = ref([]);
const loadingRecords = ref(false);
const openSessions = ref([]);
const loadingOpenSessions = ref(false);
const closeConfirm = ref({ show: false, session: null });
const usersNameFilter = ref("");
const deleteUserConfirm = ref({ show: false, user: null, saving: false });

const filters = ref({
  date_from: todayISO(),
  date_to: todayISO(),
  project: "",
  type: "",
  name: "",
  search: "",
});

function todayISO() {
  const d = new Date();
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

// Convert a YYYY-MM-DD string to a UTC ISO timestamp representing midnight in
// the browser's local timezone. The backend receives an unambiguous UTC value
// regardless of where the admin is (Mexico, Guatemala, USA, etc.).
function localDateToUTCStart(dateStr) {
  const [y, m, d] = dateStr.split("-").map(Number);
  return new Date(y, m - 1, d).toISOString();
}
// Exclusive end: midnight of the *next* day in local timezone → UTC.
function localDateToUTCEnd(dateStr) {
  const [y, m, d] = dateStr.split("-").map(Number);
  return new Date(y, m - 1, d + 1).toISOString();
}

function setToday() {
  filters.value.date_from = todayISO();
  filters.value.date_to = todayISO();
  loadAll();
}
function clearDate() {
  filters.value.date_from = "";
  filters.value.date_to = "";
  loadAll();
}

function clearFilters() {
  filters.value = {
    date_from: todayISO(),
    date_to: todayISO(),
    project: "",
    type: "",
    name: "",
    search: "",
  };
  recordsPage.value = 1;
  loadAll();
}

let searchTimer = null;
function debouncedLoad() {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    recordsPage.value = 1;
    loadRecords();
  }, 400);
}

let usersSearchTimer = null;
function debouncedLoadUsers() {
  clearTimeout(usersSearchTimer);
  usersSearchTimer = setTimeout(() => {
    usersPage.value = 1;
    loadUsers();
  }, 300);
}

async function loadAll() {
  await Promise.all([
    loadStats(),
    loadRecords(),
    loadUsers(),
    loadProjects(),
    loadOpenSessions(),
  ]);
}

async function loadStats() {
  try {
    const params = {};
    if (filters.value.date_from)
      params.date = localDateToUTCStart(filters.value.date_from);
    const { data } = await api.get("/admin/stats", { params });
    stats.value = data;
  } catch {}
}

async function loadOpenSessions() {
  loadingOpenSessions.value = true;
  try {
    const { data } = await api.get("/admin/open-sessions");
    openSessions.value = data;
  } catch {
  } finally {
    loadingOpenSessions.value = false;
  }
}

async function handleCloseSession() {
  const session = closeConfirm.value.session;
  if (!session) return;
  try {
    await api.post(`/admin/open-sessions/${session.record_id}/close`);
    openSessions.value = openSessions.value.filter(
      (s) => s.record_id !== session.record_id,
    );
  } catch {
  } finally {
    closeConfirm.value = { show: false, session: null };
  }
}

async function handleDeleteUser() {
  const user = deleteUserConfirm.value.user;
  if (!user?.id || deleteUserConfirm.value.saving) return;

  deleteUserConfirm.value.saving = true;
  try {
    await api.delete(`/admin/users/${user.id}`);
    users.value = users.value.filter((u) => u.id !== user.id);
    deleteUserConfirm.value = { show: false, user: null, saving: false };
  } catch {
    deleteUserConfirm.value.saving = false;
  }
}

function formatDuration(minutes) {
  if (minutes <= 0) return "0m";
  if (minutes < 60) return `${minutes}m`;
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  return m > 0 ? `${h}h ${m}m` : `${h}h`;
}

async function loadRecords() {
  loadingRecords.value = true;
  try {
    const params = {};
    if (filters.value.date_from)
      params.date_from = localDateToUTCStart(filters.value.date_from);
    if (filters.value.date_to)
      params.date_to = localDateToUTCEnd(filters.value.date_to);
    if (filters.value.project) params.project = filters.value.project;
    if (filters.value.type) params.type = filters.value.type;
    if (filters.value.name) params.name = filters.value.name;
    if (filters.value.search) params.search = filters.value.search;
    const { data } = await api.get("/admin/records", { params });
    records.value = data;
    recordsPage.value = 1;
  } catch {
  } finally {
    loadingRecords.value = false;
  }
}

// ─── Desviación (calculada en el navegador para respetar zona horaria local) ──
function computeDeviation(r) {
  const schedule = r.type === "entry" ? r.entry_schedule : r.exit_schedule;
  if (!schedule || schedule.length < 4) return null;
  const ts = new Date(r.timestamp);
  if (isNaN(ts.getTime())) return null;
  const actualMin = ts.getHours() * 60 + ts.getMinutes();
  const [sh, sm] = schedule.split(":").map(Number);
  const scheduledMin = sh * 60 + sm;
  return actualMin - scheduledMin;
}

function formatDeviation(dev, type) {
  if (dev === null) return "—";
  const abs = Math.abs(dev);
  const h = Math.floor(abs / 60);
  const m = abs % 60;
  const hhmm = h > 0 ? `${h}h ${m}m` : `${m}m`;
  if (dev === 0) return "Puntual";
  if (type === "entry")
    return dev > 0 ? `+${hhmm} de retraso` : `-${hhmm} anticipado`;
  return dev > 0 ? `+${hhmm} de más` : `-${hhmm} antes de tiempo`;
}

// ─── Excel export ─────────────────────────────────────────────────────────────
function exportExcel() {
  const rows = records.value.map((r) => {
    const ts = new Date(r.timestamp);
    const fecha = isNaN(ts)
      ? ""
      : ts.toLocaleDateString("es-MX", {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
        });
    const hora = isNaN(ts)
      ? ""
      : ts.toLocaleTimeString("es-MX", {
          hour: "2-digit",
          minute: "2-digit",
          hour12: false,
        });
    const horarioProg =
      r.type === "entry" ? r.entry_schedule || "—" : r.exit_schedule || "—";
    const dev = computeDeviation(r);
    return {
      Nombre: `${r.first_name || ""} ${r.last_name || ""}`.trim(),
      Correo: r.email || "",
      Proyecto: r.project_name || "",
      Tipo: r.type === "entry" ? "Entrada" : "Salida",
      Fecha: fecha,
      "Hora real": hora,
      "Horario programado": horarioProg,
      "Variación (min)": dev !== null ? dev : "",
      Variación: formatDeviation(dev, r.type),
      "Comentario del día": r.day_comment || "",
    };
  });
  const ws = XLSX.utils.json_to_sheet(rows);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, "Registros");
  const date = filters.value.date_from || todayISO();
  XLSX.writeFile(wb, `registros_${date}.xlsx`);
}

// ─── Inline comment editing ──────────────────────────────────────────────────
const editingComment = ref(null); // { user_id, record_date, text }

async function saveComment() {
  if (!editingComment.value) return;
  const { user_id, record_date, text } = editingComment.value;
  if (!text.trim()) return;
  try {
    await api.post("/admin/records/comment", {
      user_id,
      record_date,
      comment: text.trim(),
    });
    records.value.forEach((r) => {
      if (r.user_id === user_id && r.timestamp?.slice(0, 10) === record_date) {
        r.day_comment = text.trim();
      }
    });
  } catch {}
  editingComment.value = null;
}

function openCommentEditor(r) {
  const date = r.timestamp?.slice(0, 10) || todayISO();
  editingComment.value = {
    user_id: r.user_id,
    record_date: date,
    text: r.day_comment || "",
  };
}

// ─── Password reset modal ────────────────────────────────────────────────────
const pwdModal = ref({
  show: false,
  user: null,
  password: "",
  confirm: "",
  showPwd: false,
  saving: false,
  msg: { text: "", ok: false },
});

function openPwdModal(user) {
  pwdModal.value = {
    show: true,
    user,
    password: "",
    confirm: "",
    showPwd: false,
    saving: false,
    msg: { text: "", ok: false },
  };
}

async function savePwd() {
  if (
    pwdModal.value.password.length < 8 ||
    pwdModal.value.password !== pwdModal.value.confirm
  )
    return;
  pwdModal.value.saving = true;
  pwdModal.value.msg = { text: "", ok: false };
  try {
    await api.put(`/admin/users/${pwdModal.value.user.id}/password`, {
      password: pwdModal.value.password,
    });
    pwdModal.value.msg = {
      text: "Contraseña actualizada correctamente.",
      ok: true,
    };
    pwdModal.value.password = "";
    pwdModal.value.confirm = "";
    setTimeout(() => {
      pwdModal.value.show = false;
    }, 1500);
  } catch (e) {
    pwdModal.value.msg = {
      text: e.response?.data?.error || "Error al cambiar la contraseña.",
      ok: false,
    };
  } finally {
    pwdModal.value.saving = false;
  }
}

// ─── Work schedule modal ─────────────────────────────────────────────────────
const scheduleModal = ref({
  show: false,
  user: null,
  entry_time: "",
  exit_time: "",
  sat_entry_time: "",
  sat_exit_time: "",
  saving: false,
});

async function openScheduleModal(user) {
  scheduleModal.value = {
    show: true,
    user,
    entry_time: "",
    exit_time: "",
    sat_entry_time: "",
    sat_exit_time: "",
    saving: false,
  };
  try {
    const { data } = await api.get(`/admin/users/${user.id}/schedule`);
    scheduleModal.value.entry_time = data.entry_time || "";
    scheduleModal.value.exit_time = data.exit_time || "";
    scheduleModal.value.sat_entry_time = data.sat_entry_time || "";
    scheduleModal.value.sat_exit_time = data.sat_exit_time || "";
  } catch {}
}

async function saveSchedule() {
  const { user, entry_time, exit_time, sat_entry_time, sat_exit_time } =
    scheduleModal.value;
  if (!entry_time || !exit_time) return;
  scheduleModal.value.saving = true;
  try {
    await api.post(`/admin/users/${user.id}/schedule`, {
      entry_time,
      exit_time,
      sat_entry_time,
      sat_exit_time,
    });
    scheduleModal.value.show = false;
  } catch {}
  scheduleModal.value.saving = false;
}

async function loadUsers() {
  try {
    const params = {};
    if (filters.value.project) params.project = filters.value.project;
    if (usersNameFilter.value.trim()) params.name = usersNameFilter.value.trim();
    const { data } = await api.get("/admin/users", { params });
    users.value = data;
    usersPage.value = 1;
  } catch {}
}

async function loadProjects() {
  try {
    const { data } = await api.get("/admin/projects");
    projects.value = data;
  } catch {}
}

// ─── Edit user modal ─────────────────────────────────────────────────────────
const editUserModal = ref({
  show: false,
  user: null,
  first_name: "",
  last_name: "",
  project_name: "",
  email: "",
  saving: false,
  msg: { text: "", ok: false },
});

function openEditUserModal(user) {
  editUserModal.value = {
    show: true,
    user,
    first_name: user.first_name,
    last_name: user.last_name,
    project_name: user.project_name,
    email: user.email,
    saving: false,
    msg: { text: "", ok: false },
  };
}

async function saveEditUser() {
  const { user, first_name, last_name, project_name, email } =
    editUserModal.value;
  if (!first_name || !last_name || !project_name || !email) return;
  editUserModal.value.saving = true;
  editUserModal.value.msg = { text: "", ok: false };
  try {
    const { data } = await api.put(`/admin/users/${user.id}`, {
      first_name,
      last_name,
      project_name,
      email,
    });
    const u = users.value.find((u) => u.id === user.id);
    if (u) {
      u.first_name = data.first_name;
      u.last_name = data.last_name;
      u.project_name = data.project_name;
      u.email = data.email;
    }
    editUserModal.value.msg = {
      text: "Usuario actualizado correctamente.",
      ok: true,
    };
    setTimeout(() => {
      editUserModal.value.show = false;
    }, 1500);
  } catch (e) {
    editUserModal.value.msg = {
      text: e.response?.data?.error || "Error al actualizar.",
      ok: false,
    };
  } finally {
    editUserModal.value.saving = false;
  }
}

// ─── Projects management modal ────────────────────────────────────────────────
const projectsModal = ref({
  show: false,
  newName: "",
  saving: false,
  error: "",
  editing: null,   // project name currently being renamed
  editName: "",    // draft value while renaming
});

function openProjectsModal() {
  projectsModal.value = { show: true, newName: "", saving: false, error: "", editing: null, editName: "" };
}

function startEditProject(name) {
  projectsModal.value.editing = name;
  projectsModal.value.editName = name;
  projectsModal.value.error = "";
}

async function saveRenameProject(oldName) {
  const newName = projectsModal.value.editName.trim().toUpperCase();
  if (!newName || newName === oldName) {
    projectsModal.value.editing = null;
    return;
  }
  projectsModal.value.saving = true;
  projectsModal.value.error = "";
  try {
    await api.put(`/admin/projects/${encodeURIComponent(oldName)}`, { name: newName });
    const idx = projects.value.indexOf(oldName);
    if (idx !== -1) projects.value.splice(idx, 1, newName);
    projects.value = [...projects.value].sort();
    projectsModal.value.editing = null;
  } catch (e) {
    projectsModal.value.error = e.response?.data?.error || "Error al renombrar proyecto.";
  } finally {
    projectsModal.value.saving = false;
  }
}

async function saveProject() {
  const name = projectsModal.value.newName.trim().toUpperCase();
  if (!name) return;
  projectsModal.value.saving = true;
  projectsModal.value.error = "";
  try {
    await api.post("/admin/projects", { name });
    projects.value = [...projects.value, name].sort();
    projectsModal.value.newName = "";
  } catch (e) {
    projectsModal.value.error =
      e.response?.data?.error || "Error al guardar proyecto.";
  } finally {
    projectsModal.value.saving = false;
  }
}

async function deleteProject(name) {
  try {
    await api.delete(`/admin/projects/${encodeURIComponent(name)}`);
    projects.value = projects.value.filter((p) => p !== name);
  } catch {}
}

// ─── Missing records ──────────────────────────────────────────────────────────
const missingFilters = ref({
  date_from: todayISO(),
  date_to: todayISO(),
  missing_type: "both",
});
const missingRecords = ref({ missing_entry: [], missing_exit: [] });
const loadingMissing = ref(false);

async function loadMissingRecords() {
  loadingMissing.value = true;
  try {
    const params = { missing_type: missingFilters.value.missing_type };
    if (missingFilters.value.date_from)
      params.date_from = localDateToUTCStart(missingFilters.value.date_from);
    if (missingFilters.value.date_to)
      params.date_to = localDateToUTCEnd(missingFilters.value.date_to);
    const { data } = await api.get("/admin/missing-records", { params });
    missingRecords.value = data;
  } catch {
  } finally {
    loadingMissing.value = false;
  }
}

function setMissingToday() {
  missingFilters.value.date_from = todayISO();
  missingFilters.value.date_to = todayISO();
  loadMissingRecords();
}

function setMissingWeek() {
  const today = new Date();
  const dow = today.getDay();
  const monday = new Date(today);
  monday.setDate(today.getDate() - (dow === 0 ? 6 : dow - 1));
  missingFilters.value.date_from = monday.toISOString().split("T")[0];
  missingFilters.value.date_to = todayISO();
  loadMissingRecords();
}

// ─── Missing records Excel download ──────────────────────────────────────────
function downloadMissingExcel() {
  const { date_from, date_to, missing_type } = missingFilters.value;
  const label = date_from === date_to ? date_from : `${date_from}_${date_to}`;
  const wb = XLSX.utils.book_new();

  const toRow = (u) => ({
    Nombre: `${toTitleCase(u.first_name)} ${toTitleCase(u.last_name)}`,
    Proyecto: (u.project_name || "").toUpperCase(),
    Correo: u.email || "",
  });

  if (missing_type !== "exit") {
    const wsEntry = XLSX.utils.json_to_sheet(
      missingRecords.value.missing_entry.length
        ? missingRecords.value.missing_entry.map(toRow)
        : [{ Nombre: "Sin registros", Proyecto: "", Correo: "" }],
    );
    wsEntry["!cols"] = [{ wch: 30 }, { wch: 20 }, { wch: 30 }];
    XLSX.utils.book_append_sheet(wb, wsEntry, "Sin Entrada");
  }

  if (missing_type !== "entry") {
    const wsExit = XLSX.utils.json_to_sheet(
      missingRecords.value.missing_exit.length
        ? missingRecords.value.missing_exit.map(toRow)
        : [{ Nombre: "Sin registros", Proyecto: "", Correo: "" }],
    );
    wsExit["!cols"] = [{ wch: 30 }, { wch: 20 }, { wch: 30 }];
    XLSX.utils.book_append_sheet(wb, wsExit, "Sin Salida");
  }

  XLSX.writeFile(wb, `sin_registros_${label}.xlsx`);
}

// ─── Name display helpers ─────────────────────────────────────────────────────
function toTitleCase(str) {
  if (!str) return "";
  return str.replace(
    /\S+/g,
    (w) => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase(),
  );
}

function userInitials(user) {
  const first = user?.first_name?.[0] ?? "";
  const last = user?.last_name?.[0] ?? "";
  const initials = `${first}${last}`.trim().toUpperCase();
  return initials || "?";
}

// Pagination logic
const recordsPage = ref(1);
const recordsPerPage = 15;
const paginatedRecords = computed(() => {
  const start = (recordsPage.value - 1) * recordsPerPage;
  return records.value.slice(start, start + recordsPerPage);
});
const totalRecordsPages = computed(
  () => Math.ceil(records.value.length / recordsPerPage) || 1,
);

const usersPage = ref(1);
const usersPerPage = 12;
const paginatedUsers = computed(() => {
  const start = (usersPage.value - 1) * usersPerPage;
  return users.value.slice(start, start + usersPerPage);
});
const totalUsersPages = computed(
  () => Math.ceil(users.value.length / usersPerPage) || 1,
);

// Route modal + Leaflet map
const mapContainer = ref(null);
let leafletMap = null;
const routeModal = ref({
  show: false,
  user: "",
  points: [],
  duration: "",
  address: "",
  loadingAddress: false,
});

function googleMapsUrl(points) {
  if (!points.length) return "#";
  if (points.length === 1) {
    return `https://www.google.com/maps?q=${points[0].latitude},${points[0].longitude}`;
  }
  const origin = `${points[0].latitude},${points[0].longitude}`;
  const dest = `${points[points.length - 1].latitude},${points[points.length - 1].longitude}`;
  const waypoints = points
    .slice(1, -1)
    .map((p) => `${p.latitude},${p.longitude}`)
    .join("|");
  const base = `https://www.google.com/maps/dir/?api=1&origin=${origin}&destination=${dest}&travelmode=walking`;
  return waypoints
    ? `${base}&waypoints=${encodeURIComponent(waypoints)}`
    : base;
}

function destroyMap() {
  if (leafletMap) {
    leafletMap.remove();
    leafletMap = null;
  }
}

watch(
  () => routeModal.value.show,
  async (show) => {
    if (!show) {
      destroyMap();
      return;
    }
    if (!routeModal.value.points.length) return;
    await nextTick();
    if (!mapContainer.value) return;
    destroyMap();
    const points = routeModal.value.points;
    const coords = points.map((p) => [p.latitude, p.longitude]);

    // Dynamic tile layer based on theme
    leafletMap = L.map(mapContainer.value, { zoomControl: true }).setView(
      coords[0],
      17,
    );

    const tileUrl = theme.isDark
      ? "https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
      : "https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png";

    L.tileLayer(tileUrl, {
      attribution:
        '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/">CARTO</a>',
      subdomains: "abcd",
      maxZoom: 20,
    }).addTo(leafletMap);

    // Route line — glowing effect
    L.polyline(coords, { color: "#6366f1", weight: 5, opacity: 0.9 }).addTo(
      leafletMap,
    );
    L.polyline(coords, { color: "#a78bfa", weight: 2, opacity: 0.4 }).addTo(
      leafletMap,
    );

    // Start marker (green)
    const startIcon = L.divIcon({
      html: `<div style="width:18px;height:18px;background:#10b981;border:3px solid white;border-radius:50%;box-shadow:0 0 8px rgba(16,185,129,0.8);"></div>`,
      className: "",
      iconAnchor: [9, 9],
    });
    L.marker(coords[0], { icon: startIcon })
      .bindPopup(
        `<b style="color:#10b981">🟢 Inicio</b><br><span style="font-size:11px">${formatDate(points[0].recorded_at)}</span>`,
      )
      .addTo(leafletMap);

    // End marker (red)
    if (coords.length > 1) {
      const endIcon = L.divIcon({
        html: `<div style="width:18px;height:18px;background:#f43f5e;border:3px solid white;border-radius:50%;box-shadow:0 0 8px rgba(244,63,94,0.8);"></div>`,
        className: "",
        iconAnchor: [9, 9],
      });
      L.marker(coords[coords.length - 1], { icon: endIcon })
        .bindPopup(
          `<b style="color:#f43f5e">🔴 Fin</b><br><span style="font-size:11px">${formatDate(points[points.length - 1].recorded_at)}</span>`,
        )
        .addTo(leafletMap);
    }

    leafletMap.fitBounds(L.polyline(coords).getBounds(), { padding: [32, 32] });

    // Reverse geocoding with Nominatim
    routeModal.value.loadingAddress = true;
    routeModal.value.address = "";
    try {
      const { latitude: lat, longitude: lon } = points[0];
      const res = await fetch(
        `https://nominatim.openstreetmap.org/reverse?lat=${lat}&lon=${lon}&format=json&addressdetails=1`,
        { headers: { "Accept-Language": "es" } },
      );
      const geo = await res.json();
      if (geo.address) {
        const a = geo.address;
        const parts = [
          a.road,
          a.neighbourhood,
          a.suburb,
          a.city || a.town || a.village,
          a.state,
          a.postcode,
          a.country,
        ].filter(Boolean);
        // Only use unique parts to avoid duplicates like "León, León"
        routeModal.value.address =
          [...new Set(parts)].join(", ") || geo.display_name || "";
      } else {
        routeModal.value.address = geo.display_name || "";
      }
    } catch {
      routeModal.value.address = "";
    } finally {
      routeModal.value.loadingAddress = false;
    }
  },
);

async function openRouteModal(record) {
  try {
    const { data } = await api.get(`/admin/records/${record.record_id}/route`);
    // Calculate duration
    let duration = "";
    if (data.length >= 2) {
      const ms =
        new Date(data[data.length - 1].recorded_at) -
        new Date(data[0].recorded_at);
      const totalMin = Math.round(ms / 60000);
      duration =
        totalMin >= 60
          ? `${Math.floor(totalMin / 60)}h ${totalMin % 60}min`
          : `${totalMin} min`;
    }
    routeModal.value = {
      show: true,
      user: `${record.first_name} ${record.last_name}`,
      points: data,
      duration,
      address: "",
      loadingAddress: false,
    };
  } catch {}
}

// Suspicious modal
const suspiciousModal = ref({
  show: false,
  user: "",
  date: "",
  reason: "",
  ipCountry: "",
  ipCity: "",
});
function openSuspiciousModal(record) {
  suspiciousModal.value = {
    show: true,
    user: `${record.first_name} ${record.last_name}`,
    date: formatDate(record.timestamp),
    reason: record.suspicious_reason || "Sin detalle disponible",
    ipCountry: record.ip_country || "",
    ipCity: record.ip_city || "",
  };
}

// Photos modal
const photosModal = ref({
  show: false,
  user: "",
  date: "",
  site: "",
  selfie: "",
});
async function openPhotosModal(record) {
  photosModal.value = {
    show: true,
    user: `${record.first_name} ${record.last_name}`,
    date: formatDate(record.timestamp),
    site: "",
    selfie: "",
  };
  try {
    const { data } = await api.get(`/admin/records/${record.record_id}/photos`);
    photosModal.value.site = data.photo_site_path || "";
    photosModal.value.selfie = data.photo_selfie_path || "";
  } catch {}
}

function formatDate(iso) {
  return new Date(iso).toLocaleString("es-MX", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function formatDateShort(iso) {
  return new Date(iso).toLocaleDateString("es-MX", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function formatTimeOnly(iso) {
  if (!iso) return "--:--";
  return new Date(iso).toLocaleTimeString("es-MX", {
    hour: "2-digit",
    minute: "2-digit",
  });
}

const showLogoutModal = ref(false);

// ── Admin Profile ─────────────────────────────────────────
const adminFileInput = ref(null);
const adminCoverInput = ref(null);
const adminAvatarPreview = ref("");
const adminCoverPreview = ref("");
const savingAdminAvatar = ref(false);
const savingAdminCover = ref(false);
const adminAvatarOk = ref(false);
const adminAvatarErr = ref("");

const adminEmailForm = ref({ newEmail: "" });
const savingAdminEmail = ref(false);
const adminEmailMsg = ref({ text: "", ok: false });

const adminPwdForm = ref({ current: "", newPwd: "", confirm: "" });
const savingAdminPwd = ref(false);
const adminPwdMsg = ref({ text: "", ok: false });
const showAdminCurrent = ref(false);
const showAdminNew = ref(false);

const adminPwdMismatch = computed(
  () =>
    adminPwdForm.value.confirm.length > 0 &&
    adminPwdForm.value.newPwd !== adminPwdForm.value.confirm,
);

const adminPwdStrength = computed(() => {
  const p = adminPwdForm.value.newPwd;
  let score = 0;
  if (p.length >= 8) score++;
  if (/[A-Z]/.test(p)) score++;
  if (/[0-9]/.test(p)) score++;
  if (/[^A-Za-z0-9]/.test(p)) score++;
  return [
    { level: 1, color: "#ef4444", label: "Débil" },
    { level: 2, color: "#f97316", label: "Regular" },
    { level: 3, color: "#eab308", label: "Buena" },
    { level: 4, color: "#22c55e", label: "Fuerte" },
  ][Math.max(score - 1, 0)];
});

const adminInitials = computed(() => {
  const u = auth.user;
  if (!u) return "?";
  return `${u.first_name?.[0] ?? ""}${u.last_name?.[0] ?? ""}`.toUpperCase();
});

function readFile(file, callback) {
  if (!file || !file.type.startsWith("image/")) return;
  const reader = new FileReader();
  reader.onload = (ev) => callback(ev.target.result);
  reader.readAsDataURL(file);
}

function triggerAdminAvatarPick() {
  adminFileInput.value?.click();
}
function onAdminAvatarChange(e) {
  const f = e.target.files?.[0];
  if (!f) return;
  if (f.size > 500_000) {
    adminAvatarErr.value = "Máximo 500 KB";
    return;
  }
  adminAvatarErr.value = "";
  readFile(f, (v) => {
    adminAvatarPreview.value = v;
  });
}
function cancelAdminAvatar() {
  adminAvatarPreview.value = "";
  if (adminFileInput.value) adminFileInput.value.value = "";
}

async function saveAdminAvatar() {
  if (!adminAvatarPreview.value) return;
  savingAdminAvatar.value = true;
  try {
    await api.put("/profile/avatar", { avatar_url: adminAvatarPreview.value });
    auth.user.avatar_url = adminAvatarPreview.value;
    auth.persistUser();
    adminAvatarOk.value = true;
    adminAvatarPreview.value = "";
    setTimeout(() => {
      adminAvatarOk.value = false;
    }, 3000);
  } catch (e) {
    adminAvatarErr.value = e.response?.data?.error ?? "Error al guardar";
  } finally {
    savingAdminAvatar.value = false;
  }
}

function triggerAdminCoverPick() {
  adminCoverInput.value?.click();
}
function onAdminCoverChange(e) {
  const f = e.target.files?.[0];
  if (!f) return;
  readFile(f, (v) => {
    adminCoverPreview.value = v;
  });
}
function cancelAdminCover() {
  adminCoverPreview.value = "";
  if (adminCoverInput.value) adminCoverInput.value.value = "";
}

async function saveAdminCover() {
  if (!adminCoverPreview.value) return;
  savingAdminCover.value = true;
  try {
    await api.put("/profile/cover", { avatar_url: adminCoverPreview.value });
    auth.user.cover_url = adminCoverPreview.value;
    auth.persistUser();
    adminCoverPreview.value = "";
  } catch (e) {
    adminAvatarErr.value =
      e.response?.data?.error ?? "Error al guardar la portada";
  } finally {
    savingAdminCover.value = false;
  }
}

async function handleAdminEmailUpdate() {
  savingAdminEmail.value = true;
  adminEmailMsg.value = { text: "", ok: false };
  try {
    await api.put("/profile/email", { email: adminEmailForm.value.newEmail });
    auth.user.email = adminEmailForm.value.newEmail;
    auth.persistUser();
    adminEmailForm.value.newEmail = "";
    adminEmailMsg.value = {
      text: "Correo actualizado correctamente",
      ok: true,
    };
    setTimeout(() => {
      adminEmailMsg.value.text = "";
    }, 4000);
  } catch (e) {
    adminEmailMsg.value = {
      text: e.response?.data?.error ?? "Error",
      ok: false,
    };
  } finally {
    savingAdminEmail.value = false;
  }
}

async function handleAdminPwdUpdate() {
  if (adminPwdMismatch.value) return;
  savingAdminPwd.value = true;
  adminPwdMsg.value = { text: "", ok: false };
  try {
    await api.put("/profile/password", {
      current_password: adminPwdForm.value.current,
      new_password: adminPwdForm.value.newPwd,
    });
    adminPwdForm.value = { current: "", newPwd: "", confirm: "" };
    adminPwdMsg.value = { text: "Contraseña cambiada correctamente", ok: true };
    setTimeout(() => {
      adminPwdMsg.value.text = "";
    }, 4000);
  } catch (e) {
    adminPwdMsg.value = { text: e.response?.data?.error ?? "Error", ok: false };
  } finally {
    savingAdminPwd.value = false;
  }
}

const adminMemberSince = computed(() => {
  const d = auth.user?.created_at;
  if (!d) return "—";
  return new Intl.DateTimeFormat("es-MX", {
    year: "numeric",
    month: "long",
    day: "numeric",
  }).format(new Date(d));
});

async function handleLogout() {
  showLogoutModal.value = false;
  await auth.logout();
  router.push("/login");
}

onMounted(() => loadAll());

watch(activeTab, (tab) => {
  if (tab === "missing") loadMissingRecords();
});
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
.safe-bottom {
  padding-bottom: max(8px, env(safe-area-inset-bottom));
}
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media (min-width: 1024px) {
  .admin-layout {
    display: block;
  }
}
</style>
