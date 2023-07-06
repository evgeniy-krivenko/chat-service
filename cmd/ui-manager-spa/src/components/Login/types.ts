export interface ILoginFormElements extends HTMLFormControlsCollection {
  readonly login: HTMLInputElement;
  readonly password: HTMLInputElement;
}

export interface ILoginForm extends HTMLFormElement {
  readonly login: HTMLInputElement;
  readonly password: HTMLInputElement;
}
