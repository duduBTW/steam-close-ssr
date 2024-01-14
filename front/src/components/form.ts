import { html, css } from "lit";
import { customElement } from "lit/decorators.js";
import { LitElementWithTwind } from "../utils/twind";
import { createContext, consume, ContextProvider } from "@lit/context";

type SteamFormContext = boolean;
export const steamFormContext = createContext<SteamFormContext>(false);

@customElement("steam-form")
export class SteamForm extends LitElementWithTwind() {
  private _provider = new ContextProvider(this, {
    context: steamFormContext,
  });

  connectedCallback() {
    super.connectedCallback();

    this.querySelector("form")?.addEventListener("submit", async (e) => {
      e.preventDefault();

      const formElement = e.target as HTMLFormElement;

      const action = formElement.getAttribute("action");
      const method = formElement.getAttribute("method");

      if (!action || !method) {
        return;
      }

      const body = new FormData(formElement);
      const inputElements = formElement.querySelectorAll("input");

      inputElements.forEach((inputElement) => {
        inputElement.setAttribute("disabled", "true");
      });

      this._provider.setValue(true);
      await fetch(action, {
        method,
        body,
      })
        .then((res) => {
          if (res.redirected) {
            window.location.href = res.url;
          }

          return res.json();
        })
        .then((res) => res as { pg: true })
        .finally(() => {
          this._provider.setValue(false);

          inputElements.forEach((inputElement) => {
            inputElement.removeAttribute("disabled");
          });
        });
    });
  }

  protected render() {
    return html`<slot></slot>`;
  }
}

@customElement("steam-form-submit")
export class SteamFormSubmit extends LitElementWithTwind() {
  @consume({ context: steamFormContext, subscribe: true })
  isLoading: SteamFormContext = false;

  static styles = css`
    :host {
      width: min-content;
    }
  `;

  protected render() {
    const loadingContent = this.isLoading
      ? html`<loading-button-overlay />`
      : "";

    return html` <div class="relative">
      ${loadingContent}
      <slot></slot>
    </div>`;
  }
}
