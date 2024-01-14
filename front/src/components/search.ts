import { html } from "lit";
import { customElement, property, query } from "lit/decorators.js";
import { LitElementWithTwind } from "../utils/twind";
import SearchApi from "../data/search";

type EventDetail = {
  value: string;
};

const STEAM_SEARCH_CONTENT_ID = "steam-search-content";

@customElement("steam-search")
export class SteamSearch extends LitElementWithTwind() {
  @property({ attribute: true })
  open: boolean = false;

  @query(`#${STEAM_SEARCH_CONTENT_ID}`) _content!: HTMLDivElement;

  private timeout: NodeJS.Timeout | undefined;

  render() {
    if (this.open) {
      return html`
        <link
          href="https://cdn.jsdelivr.net/npm/remixicon@4.0.0/fonts/remixicon.css"
          rel="stylesheet"
        />

        <div>
          <steam-search-input
            @onClose=${this._handleClose}
            @onChange=${this._handleChange}
          ></steam-search-input>
          <div id=${STEAM_SEARCH_CONTENT_ID} />
        </div>
      `;
    }

    return html`<slot></slot>`;
  }

  connectedCallback(): void {
    super.connectedCallback();
    this.querySelector("#search-trigger")?.addEventListener("click", (e) => {
      e.preventDefault();
      this.open = !this.open;
    });
  }

  private _handleClose() {
    this.open = false;
    this._clearContent();
  }

  private _clearContent() {
    this._content.innerHTML = "";
  }

  private async _fetchFastSearchList(value: string) {
    if (!value) {
      this._clearContent();
      return;
    }

    this._content.innerHTML = await SearchApi.fetchFastSearchList(value);
  }

  private _handleChange(event: CustomEvent<EventDetail>) {
    if (this.timeout) {
      clearTimeout(this.timeout);
    }

    this.timeout = setTimeout(
      () => this._fetchFastSearchList(event.detail.value),
      300
    );
  }
}

@customElement("steam-search-input")
export class SteamSearchInput extends LitElementWithTwind() {
  render() {
    return html`
      <link
        href="https://cdn.jsdelivr.net/npm/remixicon@4.0.0/fonts/remixicon.css"
        rel="stylesheet"
      />

      <div
        href="/search"
        id="search-trigger"
        class="h-9 flex items-center justify-center px-1 rounded-t bg-gray-800 relative w-80"
        id="dropdown-trigger"
      >
        <i
          class="ri-search-line absolute top-2/4 -translate-y-2/4 left-3 text-gray-500"
        ></i>

        <input
          placeholder="Search..."
          class="flex bg-transparent h-7 w-full rounded-md border border-gray-700 px-2 py-1 text-xs shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 pl-7"
          type="text"
          @input=${this._handleChange}
          autofocus
        />

        <button @click=${this._handleClose}>
          <i class="ri-close-fill px-2 ml-1"></i>
        </button>
      </div>
    `;
  }

  private _handleClose() {
    this.dispatchEvent(
      new CustomEvent("onClose", {
        bubbles: true,
        composed: true,
      })
    );
  }

  private _handleChange(e: Event) {
    const value = (e.target as HTMLInputElement).value;

    this.dispatchEvent(
      new CustomEvent<EventDetail>("onChange", {
        bubbles: true,
        composed: true,
        detail: {
          value,
        },
      })
    );
  }
}

@customElement("steam-search-game-list")
export class SteamSearchGamesList extends LitElementWithTwind() {
  @property({ attribute: true })
  value: string = "";

  render() {
    if (!this.value) {
      return "";
    }

    return html``;
  }
}
