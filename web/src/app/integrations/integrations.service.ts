import {Inject, Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {map} from 'rxjs/internal/operators';
import {Observable} from 'rxjs';
import {Integration} from './integration.model';
import {AuthService} from '../services/auth.service';

@Injectable()
export class IntegrationsService {
  constructor(@Inject('api') private api: string,
              private http: HttpClient,
              private auth: AuthService) {
  }

  list(): Observable<Integration[]> {
    return this.http.get(`${this.api}integrations`, {headers: this.headers()}).pipe(
      map((response: Integration[]) => {
        const data: Integration[] = [];
        response.map((item: Integration) => {
          data.push(
            new Integration(
              item.id,
              item.title,
              item.token,
              item.channel,
              item.secret,
              item.proxyAddress,
              item.proxyUser,
              item.proxyPass,
            ));
        });

        return data;
      })
    );
  }

  save(model: Integration): Observable<any> {
    let url = `${this.api }integrations`;
    if (model.id) {
      url = `${url}/${model.id}`;
    }

    return this.http.post(url, JSON.stringify(model), {headers: this.headers()});
  }

  remove(model: Integration): Observable<any> {
    const url = `${this.api }integrations/${model.id}`;
    return this.http.delete(url, {headers: this.headers()});
  }

  headers(): HttpHeaders {
    let headers = new HttpHeaders();
    headers = headers.append('Content-Type', 'application/json');
    return headers.append('Authorization', 'Bearer ' + this.auth.getToken());
  }
}
