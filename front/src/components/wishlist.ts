import { html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { LitElementWithTwind } from "../utils/twind";

const wishlistActionNodeType = "wishlist-action";
function wishListGame() {
  return new Promise<string>((resolve, reject) =>
    fetch("/api/game/wishlist")
      .then((r) => {
        if (r.status !== 200) {
          throw new Error(r.statusText);
        }

        return r.text();
      })
      .then(resolve)
      .catch(reject)
  );
}

@customElement(wishlistActionNodeType)
export class WishlistAction extends LitElementWithTwind() {
  @property({ type: String })
  gameId?: string;

  @property({ type: Boolean })
  isLoading = false;

  constructor() {
    super();
    this.addEventListener("click", this._wishListGame);
  }

  protected render() {
    const loadingContent = this.isLoading ? html`<loading-content />` : "";

    return html` <div class="relative">
      ${loadingContent}
      <slot></slot>
    </div>`;
  }

  private async _wishListGame(e: Event) {
    e.preventDefault();

    if (this.isLoading) {
      return;
    }

    const buttonElement = this.querySelector("button");
    if (!buttonElement) {
      return;
    }

    this.isLoading = true;
    buttonElement.setAttribute("disabled", "true");

    try {
      this.innerHTML = await wishListGame();
    } finally {
      buttonElement?.removeAttribute("disabled");
      this.isLoading = false;
    }
  }
}

@customElement("loading-content")
export class LoadingContent extends LitElementWithTwind() {
  protected render() {
    return html` <link
        href="https://cdn.jsdelivr.net/npm/remixicon@4.0.0/fonts/remixicon.css"
        rel="stylesheet"
      />

      <div class="absolute w-full h-full bg-zinc-900/60 backdrop-blur rounded grid place-items-center z-20">
        <i class="ri-loader-fill animate-spin"></i>
      </divc>`;
  }
}
