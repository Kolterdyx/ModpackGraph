import { FormControl } from '@angular/forms';

export type Form<T> = {
  [K in keyof T]: FormControl<T[K] | null>;
};

export type Nullable<T> = {
  [K in keyof T]: T[K] | null;
}
