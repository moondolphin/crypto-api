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

  // =========================
  // PASO 3: Última cotización
  // =========================
  const btnLastPrice = $("btn-last-price");
  if (btnLastPrice) {
    btnLastPrice.addEventListener("click", async () => {
      const symbol = toUpper($("lp-symbol")?.value || "");
      const provider = $("lp-provider")?.value || "";
      const currency = $("lp-currency")?.value || "";

      if (!symbol) {
        alert("Ingresá un symbol (BTC, ETH, etc)");
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
        $("lp-out-ts").textContent = pick(data, "timestamp", "Timestamp") ?? "-";
      } catch (err) {
        console.error(err);
        alert("Error consultando precio (mirá Console/Network)");
      }
    });
  }

  // =========================
  // PASO 4: Tabla de quotes
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
        const quotedAt =
          pick(q, "quoted_at", "quotedAt", "QuotedAt", "Quoted_At") ?? "-";

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
    // backend real: snake_case
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
});
