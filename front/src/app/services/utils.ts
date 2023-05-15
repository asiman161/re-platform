import { environment } from '../../environments/environment';

export function makePathPrefix(): string {
  return `${environment.schema}${environment.host}:${environment.port}`
}
