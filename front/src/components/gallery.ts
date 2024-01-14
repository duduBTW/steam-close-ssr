import { html, css } from "lit";
import { customElement, property } from "lit/decorators.js";
import { LitElementWithTwind } from "../utils/twind";
import { createContext, consume, ContextProvider } from "@lit/context";
import { styleMap } from "lit/directives/style-map.js";

type SteamGalleryContext = string | null;
export const steamGalleryContext = createContext<SteamGalleryContext>(null);

@customElement("steam-gallery")
export class SteamGallery extends LitElementWithTwind() {
  private _provider = new ContextProvider(this, {
    context: steamGalleryContext,
  });

  connectedCallback() {
    super.connectedCallback();
    let prevItem: string | null = null;

    this.querySelectorAll("steam-gallery-sub")?.forEach((item) =>
      item.addEventListener("mouseenter", (e) => {
        this._provider.setValue(
          (e.target as HTMLDivElement).getAttribute("imageUrl")
        );
      })
    );

    this.querySelectorAll("steam-gallery-sub").forEach((item) =>
      item.addEventListener("mouseleave", async () => {
        prevItem = this._provider.value;

        setTimeout(() => {
          if (prevItem !== this._provider.value) {
            return;
          }

          this._provider.setValue(null);
        }, 10);
      })
    );
  }

  protected render() {
    return html`<slot></slot>`;
  }
}

@customElement("steam-gallery-main")
export class SteamGalleryMain extends LitElementWithTwind() {
  protected render() {
    return html`<div class="relative rounded-lg overflow-hidden">
      <slot></slot>
      <steam-gallery-main-preview />
    </div>`;
  }
}

@customElement("steam-gallery-main-preview")
export class SteamGalleryMainPreview extends LitElementWithTwind() {
  @consume({ context: steamGalleryContext, subscribe: true })
  selectedImage: SteamGalleryContext = null;

  static styles = css`
    img {
      animation: fadeIn 0.16s cubic-bezier(0.4, 0, 0.2, 1);
    }

    @keyframes fadeIn {
      from {
        opacity: 0;
      }

      to {
        opacity: 1;
      }
    }
  `;

  protected render() {
    if (!this.selectedImage) {
      return;
    }

    return html`<img
      class="absolute w-full h-full top-0 left-0 object-cover"
      src=${this.selectedImage}
    />`;
  }
}

@customElement("steam-gallery-sub")
export class SteamGallerySub extends LitElementWithTwind() {
  @consume({ context: steamGalleryContext, subscribe: true })
  selectedImage: SteamGalleryContext = null;

  @property({ attribute: true })
  imageUrl: string = "";

  protected render() {
    const styles = {
      opacity:
        this.selectedImage && this.imageUrl !== this.selectedImage
          ? "0.2"
          : "1",
    };

    return html`<div
      style=${styleMap(styles)}
      data-image-url=${this.imageUrl}
      id="steam-gallery-item"
      class="transition-opacity"
    >
      <slot></slot>
    </div>`;
  }
}
