console.log("✅ app.js cargó");

document.addEventListener("DOMContentLoaded", () => {
 
  // =========================
  // Helpers
  // =========================
 
  const $ = (id) => document.getElementById(id);

  const pick = (obj, ...keys) => {
    for (const k of keys) {
      if (obj && obj[k] !== undefined && obj[k] !== null) return obj[k];
    }
    return undefined;
  };

  const toUpper = (s) => (s || "").trim().toUpperCase();

function showToast(message, variant = "dark") {
    const toastEl = document.getElementById("app-toast");
    const toastBody = document.getElementById("app-toast-body");
    if (!toastEl || !toastBody || !window.bootstrap) return;

    toastBody.textContent = message;

    // variants: success, danger, warning, info, dark
    toastEl.classList.remove("text-bg-success", "text-bg-danger", "text-bg-warning", "text-bg-info", "text-bg-dark");
    toastEl.classList.add(`text-bg-${variant}`);

    const toast = window.bootstrap.Toast.getOrCreateInstance(toastEl, {
      delay: 3500,
      autohide: true,
    });
    toast.show();
  }

function showModal(message, title = "Atención", icon = "ℹ️") {
    const modalEl = document.getElementById("appMessageModal");
    const bodyEl = document.getElementById("appMessageBody");
    const titleEl = document.getElementById("appMessageTitle");
    const iconEl = document.getElementById("appMessageIcon");

    if (!modalEl || !bodyEl || !titleEl || !iconEl || !window.bootstrap) return;

    titleEl.textContent = title;
    bodyEl.textContent = message;
    iconEl.textContent = icon;

    const modal = window.bootstrap.Modal.getOrCreateInstance(modalEl, {
      backdrop: "static",
      keyboard: true,
    });

    modal.show();
  }

  function formatLocalDate(isoString) {
    if (!isoString) return "-";

    const d = new Date(isoString);
    if (isNaN(d.getTime())) return isoString;

    return new Intl.DateTimeFormat(undefined, {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    }).format(d);
  }



  // =========================
  // Última cotización
  // =========================
  const btnLastPrice = $("btn-last-price");
  if (btnLastPrice) {
    btnLastPrice.addEventListener("click", async () => {
      const symbol = toUpper($("lp-symbol")?.value || "");
      const provider = $("lp-provider")?.value || "";
      const currency = $("lp-currency")?.value || "";

    if (!symbol) {
        showModal( "Ingresá un symbol válido (BTC, ETH, UNI, etc).",
            "Dato requerido","⚠️" );
        return;
    }


      const params = new URLSearchParams({ symbol, provider, currency });

      try {
        const res = await fetch(`/api/v1/crypto/price?${params.toString()}`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();

        $("lp-out-symbol").textContent = pick(data, "symbol", "Symbol") ?? "-";
        $("lp-out-price").textContent = pick(data, "price", "Price") ?? "-";
        $("lp-out-provider").textContent = pick(data, "provider", "Provider") ?? "-";
        const tsRaw = pick(data, "timestamp", "Timestamp");
        $("lp-out-ts").textContent = formatLocalDate(tsRaw);
      } catch (err) {
        console.error(err);
        showModal( "Error consultando precio", "Atención", "❌");
      }
    });
  }

  // =========================
  // Tabla de quotes
  // =========================
  const tbody = $("quotes-tbody");
  const summaryEl = $("quotes-summary");

  const btnRefresh = $("btn-refresh-table");
  const btnApply = $("btn-apply-filters");
  const btnClear = $("btn-clear-filters");
  const btnPrev = $("btn-prev");
  const btnNext = $("btn-next");

  const fSymbol = $("f-symbol");
  const fProvider = $("f-provider");
  const fCurrency = $("f-currency");
  const fFrom = $("f-from");
  const fTo = $("f-to");
  const fMin = $("f-min");
  const fMax = $("f-max");
  const fPage = $("f-page");
  const fPageSize = $("f-page-size");

  function setLoading() {
    tbody.innerHTML =
      '<tr><td colspan="5" class="text-center text-muted py-4">Cargando...</td></tr>';
  }

  function setEmpty() {
    tbody.innerHTML =
      '<tr><td colspan="5" class="text-center text-muted py-4">Sin resultados.</td></tr>';
  }

  function escapeHtml(s) {
    return String(s)
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }

  function renderRows(items) {
    tbody.innerHTML = items
      .map((q) => {
        const symbol = pick(q, "symbol", "Symbol") ?? "-";
        const provider = pick(q, "provider", "Provider") ?? "-";
        const currency = pick(q, "currency", "Currency") ?? "-";
        const price = pick(q, "price", "Price") ?? "-";
        const quotedAtRaw =
            pick(q, "quoted_at", "quotedAt", "QuotedAt", "Quoted_At") ?? "-";

        const quotedAt = formatLocalDate(quotedAtRaw);


        return (
          "<tr>" +
          `<td class="fw-semibold">${escapeHtml(symbol)}</td>` +
          `<td>${escapeHtml(provider)}</td>` +
          `<td>${escapeHtml(currency)}</td>` +
          `<td class="text-end">${escapeHtml(price)}</td>` +
          `<td>${escapeHtml(quotedAt)}</td>` +
          "</tr>"
        );
      })
      .join("");
  }

  function updateSummary(summary) {
    const totalItems =
      pick(summary, "total_items", "totalItems", "TotalItems") ?? "-";
    const totalPages =
      pick(summary, "total_pages", "totalPages", "TotalPages") ?? "-";
    const page = pick(summary, "page", "Page") ?? "-";
    const pageSize =
      pick(summary, "page_size", "pageSize", "PageSize") ?? "-";

    summaryEl.textContent = `Total items: ${totalItems} | Total pages: ${totalPages} | Page: ${page} | Page size: ${pageSize}`;

    const p = Number(page);
    const tp = Number(totalPages);
    btnPrev.disabled = !Number.isFinite(p) || p <= 1;
    btnNext.disabled =
      Number.isFinite(tp) && Number.isFinite(p) ? p >= tp : false;
  }

  function buildQuery() {
    // snake_case
    const params = new URLSearchParams();

    const symbol = toUpper(fSymbol.value);
    if (symbol) params.set("symbol", symbol);
    if (fProvider.value) params.set("provider", fProvider.value);
    if (fCurrency.value) params.set("currency", fCurrency.value);

    if (fFrom.value) params.set("from", fFrom.value);
    if (fTo.value) params.set("to", fTo.value);

    if (fMin.value) params.set("min_price", fMin.value);
    if (fMax.value) params.set("max_price", fMax.value);

    params.set("page", fPage.value || "1");
    params.set("page_size", fPageSize.value || "50");

    return params.toString();
  }

  async function loadQuotes() {
    setLoading();
    try {
      const qs = buildQuery();
      const res = await fetch(`/api/v1/quotes?${qs}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);

      const data = await res.json();
      const items = pick(data, "items", "Items") || [];
      const summary = pick(data, "summary", "Summary") || {};

      if (!items.length) setEmpty();
      else renderRows(items);

      updateSummary(summary);

      const serverPage = pick(summary, "page", "Page");
      const serverPageSize = pick(summary, "page_size", "pageSize", "PageSize");
      if (serverPage) fPage.value = String(serverPage);
      if (serverPageSize) fPageSize.value = String(serverPageSize);
    } catch (err) {
      console.error(err);
      tbody.innerHTML =
        '<tr><td colspan="5" class="text-center text-danger py-4">Error cargando quotes (mirá Console/Network)</td></tr>';
      summaryEl.textContent =
        "Total items: - | Total pages: - | Page: - | Page size: -";
    }
  }

  function clearFilters() {
    fSymbol.value = "";
    fProvider.value = "";
    fCurrency.value = "";
    fFrom.value = "";
    fTo.value = "";
    fMin.value = "";
    fMax.value = "";
    fPage.value = "1";
    fPageSize.value = "50";
  }

  // Eventos
  btnRefresh?.addEventListener("click", loadQuotes);

  btnApply?.addEventListener("click", () => {
    fPage.value = "1";
    loadQuotes();
  });

  btnClear?.addEventListener("click", () => {
    clearFilters();
    loadQuotes();
  });

  btnPrev?.addEventListener("click", () => {
    const p = Math.max(1, Number(fPage.value || "1") - 1);
    fPage.value = String(p);
    loadQuotes();
  });

  btnNext?.addEventListener("click", () => {
    const p = Number(fPage.value || "1") + 1;
    fPage.value = String(p);
    loadQuotes();
  });

  // Primera carga
  loadQuotes();

    // =========================
  // Auth (JWT)
  // =========================
  const TOKEN_KEY = "crypto_api_token";

  const getToken = () => localStorage.getItem(TOKEN_KEY) || "";
  const setToken = (t) => localStorage.setItem(TOKEN_KEY, t);
  const clearToken = () => localStorage.removeItem(TOKEN_KEY);
  const isLoggedIn = () => Boolean(getToken());

  const btnOpenLogin = $("btn-open-login");
  const btnOpenRegister = $("btn-open-register"); 
  const btnOpenCoins = $("btn-open-coins");
  const btnRunRefresh = $("btn-run-refresh");

  const loginError = $("login-error");
  const loginForm = $("login-form");
  const loginEmail = $("login-email");
  const loginPassword = $("login-password");

  function setPrivateEnabled(enabled) {
    if (btnOpenCoins) btnOpenCoins.disabled = !enabled;
    if (btnRunRefresh) btnRunRefresh.disabled = !enabled;
  }

  function setNavbarAuthState() {
    const logged = isLoggedIn();

    // Enable/disable privados
    setPrivateEnabled(logged);

    // Convertimos el botón Login en Logout cuando hay token
    if (btnOpenLogin) {
      btnOpenLogin.textContent = logged ? "Logout" : "Login";
      btnOpenLogin.classList.toggle("btn-warning", !logged);
      btnOpenLogin.classList.toggle("btn-outline-light", logged);
    }

    if (btnOpenRegister) {
      btnOpenRegister.disabled = true;
    }
  }

  async function authFetch(url, options = {}) {
    const token = getToken();
    const headers = new Headers(options.headers || {});
    if (token) headers.set("Authorization", `Bearer ${token}`);
    return fetch(url, { ...options, headers });
  }

  function showLoginError(msg) {
    if (!loginError) return;
    loginError.textContent = msg;
    loginError.classList.remove("d-none");
    showModal( "Error consultando precio");
  }

  function clearLoginError() {
    if (!loginError) return;
    loginError.textContent = "";
    loginError.classList.add("d-none");
  }

  // Inicial
  setNavbarAuthState();



  // =========================
  // Helpers
  // =========================
  // Abrir modal de login
  btnOpenLogin?.addEventListener("click", () => {
    if (isLoggedIn()) {
      clearToken();
      setNavbarAuthState();
      return;
    }

    clearLoginError();
    loginEmail.value = "";
    loginPassword.value = "";

    // Bootstrap modal
    const el = document.getElementById("loginModal");
    if (!el) return;

    const modal = window.bootstrap?.Modal?.getOrCreateInstance(el);
    modal?.show();
  });

  // Submit login
  loginForm?.addEventListener("submit", async (e) => {
    e.preventDefault();
    clearLoginError();

    const email = (loginEmail?.value || "").trim();
    const password = loginPassword?.value || "";

    if (!email || !password) {
      //showLoginError("Email y password son requeridos.");
      showModal( "Email y password son requeridos.");
      return;
    }

    try {
      const res = await fetch("/api/v1/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        // Intentamos leer mensaje si viene JSON, sino usamos status
        let msg = `Login inválido (HTTP ${res.status})`;
        try {
          const t = await res.text();
          if (t) msg = t;
        } catch (_) {}
        //showLoginError(msg);
        showModal(msg);
        return;
      }

      const data = await res.json();
      const token = data.access_token || data.accessToken || data.AccessToken;

      if (!token) {
        //showLoginError();
         showModal("El server no devolvió access_token.");
        return;
      }

      setToken(token);
      setNavbarAuthState();

      // Cerrar modal
      const el = document.getElementById("loginModal");
      const modal = el ? window.bootstrap?.Modal?.getOrCreateInstance(el) : null;
      modal?.hide();
    } catch (err) {
      console.error(err);
      //showLoginError("Error de red en login (mirá Console).");
      showModal("Error de red en login (mirá Console).");
    }
  });
    // =========================
  //  POST Refresh (manual)
  // =========================
  btnRunRefresh?.addEventListener("click", async () => {
    try {
      btnRunRefresh.disabled = true;
      const oldText = btnRunRefresh.textContent;
      btnRunRefresh.textContent = "Refreshing...";

      const res = await authFetch("/api/v1/job/refresh", {
        method: "POST",
      });

      if (!res.ok) {
        const body = await res.text().catch(() => "");
        throw new Error(body || `HTTP ${res.status}`);
      }

      const data = await res.json();

      const coinsProcessed = data.coins_processed ?? data.coinsProcessed ?? data.CoinsProcessed ?? "-";
      const quotesSaved = data.quotes_saved ?? data.quotesSaved ?? data.QuotesSaved ?? "-";
      const failed = data.failed ?? data.Failed ?? "-";

      showToast(`Refresh OK ✅  Coins: ${coinsProcessed} | Quotes: ${quotesSaved} | Failed: ${failed}`, "success");

      // Recargar tabla de quotes
      if (typeof loadQuotes === "function") {
        await loadQuotes();
      }
    } catch (err) {
      console.error(err);
      showToast(`Refresh falló ❌ ${err?.message || err}`, "danger");
    } finally {
      btnRunRefresh.disabled = false;
      btnRunRefresh.textContent = "POST Refresh (manual)";
    }
  });

    // =========================
  // Gestionar Coins (modal + acciones)
  // =========================
  const coinSymbolInput = document.getElementById("coin-symbol");
  const btnCoinDisable = document.getElementById("btn-coin-disable");
  const btnCoinEnable = document.getElementById("btn-coin-enable");

  function getCoinSymbol() {
    return (coinSymbolInput?.value || "").trim().toUpperCase();
  }

  function openCoinsModal() {
    const el = document.getElementById("coinsModal");
    if (!el || !window.bootstrap) return;

    if (coinSymbolInput) coinSymbolInput.value = "";

    const modal = window.bootstrap.Modal.getOrCreateInstance(el);
    modal.show();
  }

  // Abrir modal desde el botón del dashboard
  btnOpenCoins?.addEventListener("click", () => {
    if (!isLoggedIn()) {
      showModal("Tenés que iniciar sesión para gestionar coins.", "Acceso requerido", "🔒");
      return;
    }
    openCoinsModal();
  });

  // Alta / Habilitar moneda (POST /api/v1/coins)
  btnCoinEnable?.addEventListener("click", async () => {
    const symbol = getCoinSymbol();
    if (!symbol) {
      showModal("Ingresá un symbol (BTC, ETH, UNI...).", "Dato requerido", "⚠️");
      return;
    }

    try {
      btnCoinEnable.disabled = true;

      const res = await authFetch("/api/v1/coins", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          symbol,
          enabled: true, // default true (alta/habilitar)
          // coingecko_id / binance_symbol: opcionales
        }),
      });

      if (!res.ok) {
        const body = await res.text().catch(() => "");
        throw new Error(body || `HTTP ${res.status}`);
      }

      await res.json().catch(() => ({}));
      showToast(`Moneda habilitada ✅ ${symbol}`, "success");

    } catch (err) {
      console.error(err);
      showModal(
        `No se pudo dar de alta/habilitar la moneda.\n${err?.message || err}`,
        "Error",
        "❌"
      );
    } finally {
      btnCoinEnable.disabled = false;
    }
  });

  // Deshabilitar moneda (PUT /api/v1/coins/{symbol} con enabled=false)
  btnCoinDisable?.addEventListener("click", async () => {
    const symbol = getCoinSymbol();
    if (!symbol) {
      showModal("Ingresá un symbol (BTC, ETH, UNI...).", "Dato requerido", "⚠️");
      return;
    }

    try {
      btnCoinDisable.disabled = true;

      const res = await authFetch(`/api/v1/coins/${encodeURIComponent(symbol)}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ enabled: false }),
      });

      if (!res.ok) {
        const body = await res.text().catch(() => "");
        throw new Error(body || `HTTP ${res.status}`);
      }

      await res.json().catch(() => ({}));
      showToast(`Moneda deshabilitada ✅ ${symbol}`, "success");

    } catch (err) {
      console.error(err);
      showModal(
        `No se pudo deshabilitar la moneda.\n${err?.message || err}`,
        "Error",
        "❌"
      );
    } finally {
      btnCoinDisable.disabled = false;
    }
  });


});
