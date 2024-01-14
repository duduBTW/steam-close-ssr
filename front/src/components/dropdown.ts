import { css, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { LitElementWithTwind } from "../utils/twind";
import { classMap } from "lit/directives/class-map.js";
import { createContext, consume, ContextProvider } from "@lit/context";

type DropdownContext = boolean;
export const dropdownContext = createContext<DropdownContext>(false);

@customElement("steam-dropdown")
export class UserNavDropdownTrigger extends LitElementWithTwind() {
  private _provider = new ContextProvider(this, {
    context: dropdownContext,
  });

  protected render() {
    return html`<div>
      <slot></slot>
    </div>`;
  }

  connectedCallback() {
    super.connectedCallback();

    this.querySelector("#dropdown-trigger")?.addEventListener("click", (e) => {
      e.preventDefault();
      this._provider.setValue(!this._provider.value);
    });
  }
}

@customElement("dropdown-content")
export class UserNavDropdownContent extends LitElementWithTwind() {
  @consume({ context: dropdownContext, subscribe: true })
  open: DropdownContext = false;

  @property({ attribute: false })
  _width = "0px";

  protected render() {
    const containerClasses = {
      "animation-enabled": !!this.open,
      "animation-hidden": !this.open,
      "animation-grid": true,
    };

    const contentClasses = {
      [`w-[${this._width}]`]: !!this.open,
      [`w-0`]: !this.open,
      "animation-content": true,
    };

    return html`
      <div class=${classMap(containerClasses)}>
        ${this.open
          ? html`<div
              class="flex flex-col rounded-b bg-gray-800 overflow-hidden ${classMap(
                contentClasses
              )}"
            >
              <slot></slot>
            </div>`
          : ""}
      </div>
    `;
  }

  connectedCallback() {
    super.connectedCallback();
    this._width = window.getComputedStyle(this.firstChild as HTMLElement).width;
  }

  static styles = css`
    .animation-enabled {
      grid-template-rows: 1fr;
    }
    .animation-hidden {
      grid-template-rows: 0fr;
    }
    .animation-grid {
      display: grid;
      transition: grid-template-rows 0.2s cubic-bezier(0.4, 0, 0.2, 1),
        margin 0.2s cubic-bezier(0.4, 0, 0.2, 1),
        padding 0.2s cubic-bezier(0.4, 0, 0.2, 1), width 0.2s ease;
      overflow: hidden;
    }

    .animation-content {
      transition: width 1s ease-in-out;
    }

    :host {
      display: block !important;
    }
  `;
}
