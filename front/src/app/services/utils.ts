import { environment } from '../../environments/environment';

export function makePathPrefix(): string {
  return `${environment.schema}${window.location.hostname}:${environment.port}`
}
