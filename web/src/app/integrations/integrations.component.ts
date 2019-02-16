import {Component, Inject, OnInit} from '@angular/core';
import {IntegrationsService} from "./integrations.service";
import {AuthService} from "../services/auth.service";
import {Integration} from "./integration.model";

@Component({
  selector: 'app-integrations',
  templateUrl: './integrations.component.html',
  styleUrls: ['./integrations.component.css']
})
export class IntegrationsComponent implements OnInit {

  models: Integration[];
  loading: boolean;
  message: string;

  constructor(
    @Inject('server') private server: string,
    private storage: IntegrationsService,
    public auth: AuthService) {
    this.loading = true;
    this.message = '';
  }

  ngOnInit() {
    this.message = '';
    this.storage.list().subscribe(data => {
      this.models = data;
      this.loading = false;
    });
  }

  // add new form
  add() {
    this.models.push(new Integration());
  }

  remove(index) {
    this.message = '';
    this.loading = true;
    if (this.models[index] && this.models[index].id) {
      this.storage.remove(this.models[index])
        .subscribe(data => {
            this.loading = false;
          },
          err => {
            this.loading = false;
            if (err.error.message) {
              this.message = err.error.message.replace(/;/g, '. ');
            }
          });
    } else {
      this.loading = false;
    }
    this.models.splice(index, 1);
  }

  // save form
  save(feed: Integration) {
    this.message = '';
    this.loading = true;
    this.storage.save(feed)
      .subscribe(data => {
        this.loading = false;
      }, err => {
        this.loading = false;
        if (err.error.message) {
          this.message = err.error.message.replace(/;/g, '. ');
        }
      });
    console.log(feed);
  }

}
