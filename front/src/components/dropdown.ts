import { css, html } from "lit";
import { customElement } from "lit/decorators.js";
import { LitElementWithTwind } from "../utils/twind";
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

  protected render() {
    if (!this.open) {
      return "";
    }

    return html`<div
      class="flex flex-col rounded-b bg-gray-800 overflow-hidden"
    >
      <slot></slot>
    </div>`;
  }

  static styles = css`
    :host {
      display: block !important;
    }
  `;
}
