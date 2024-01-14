import { html, css } from "lit";
import { customElement, property } from "lit/decorators.js";
import { classMap } from "lit/directives/class-map.js";
import { LitElementWithTwind } from "../utils/twind";

const steamCarouselPaginationNodeType = "steam-carousel-pagination";
const steamCarouselItemNodeType = "steam-carousel-item";

@customElement("steam-carousel")
export class SteamCarousel extends LitElementWithTwind() {
  @property({ type: Number })
  activeSlide: number = 1;

  @property({ type: Boolean })
  loop: boolean = false;

  _numberOfPages = 0;

  protected render() {
    return html`
      <link
        href="https://cdn.jsdelivr.net/npm/remixicon@4.0.0/fonts/remixicon.css"
        rel="stylesheet"
      />

      <div class="flex justify-between bg-gray-800 p-4 gap-4 rounded-lg">
        <button @click=${this._slidePrevious}>
          <i class="ri-arrow-left-s-line"></i>
        </button>
        <div class="overflow-hidden flex-1">
          <slot
            style="transform: translateX(-${(this.activeSlide - 1) * 100}%);"
          ></slot>
        </div>
        <button @click=${this._slideNext}>
          <i class="ri-arrow-right-s-line"></i>
        </button>
      </div>

      <steam-carousel-pagination
        activeSlide=${this.activeSlide}
        numberOfPages=${this._numberOfPages}
        @onPageClick=${this._slideTo}
      />
    `;
  }

  constructor() {
    super();
    let _numberOfPages = 0;

    this.childNodes.forEach((node) => {
      if (
        node.nodeName.toLowerCase() !== steamCarouselItemNodeType.toLowerCase()
      ) {
        return;
      }

      _numberOfPages++;
    });

    this._numberOfPages = _numberOfPages;
  }

  static styles = css`
    slot {
      display: flex;
      flex-wrap: nowrap;
      transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }
  `;

  _slideTo(e: OnPageClickCustomEvent) {
    this.activeSlide = e.detail.page;
  }

  _slideNext() {
    const isLastItem = this.activeSlide === this._numberOfPages;

    if (isLastItem && !this.loop) {
      return;
    }

    if (isLastItem && this.loop) {
      this.activeSlide = 1;
      return;
    }

    this.activeSlide++;
  }

  _slidePrevious() {
    const isFirstItem = this.activeSlide === 1;

    if (isFirstItem && !this.loop) {
      return;
    }

    if (isFirstItem && this.loop) {
      this.activeSlide = this._numberOfPages;
      return;
    }

    this.activeSlide--;
  }
}

@customElement(steamCarouselItemNodeType)
export class SteamCarouselItem extends LitElementWithTwind() {
  static styles = css`
    :host {
      width: 100%;
      flex-grow: 1;
      flex-shrink: 0;
    }

    img {
      pointer-events: none;
      user-select: none;
    }
  `;

  protected render() {
    return html`
      <div class="w-full">
        <slot></slot>
      </div>
    `;
  }
}

type OnPageClickCustomEvent = CustomEvent<{
  page: number;
}>;

@customElement(steamCarouselPaginationNodeType)
export class SteamCarouselPagination extends LitElementWithTwind() {
  @property()
  numberOfPages: number = 0;

  @property({ type: Number })
  activeSlide: number = 1;

  protected render() {
    return html`
      <div class="flex justify-center gap-2">
        ${Array.from({ length: this.numberOfPages }).map((_, index) => {
          const selected = this.activeSlide === index + 1;
          const classes = {
            "bg-black": !selected,
            "bg-blue-600": selected,
          };

          return html`<button
            data-page=${index + 1}
            @click=${this._handlePageClick}
            class="w-8 h-2 rounded-full ${classMap(classes)} transition-colors"
          />`;
        })}
      </div>
    `;
  }

  _handlePageClick(e: Event) {
    this.dispatchEvent(
      new CustomEvent("onPageClick", {
        detail: {
          page: Number((e.target as Element).getAttribute("data-page")),
        },
      })
    );
  }
}
