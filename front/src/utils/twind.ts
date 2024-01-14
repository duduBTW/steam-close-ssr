import { LitElement } from "lit";
import install from "@twind/with-web-components";
import config from "./twind.config";

export const LitElementWithTwind = () => install(config)(LitElement);
