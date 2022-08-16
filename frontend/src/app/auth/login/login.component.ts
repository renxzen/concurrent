import { Component, OnInit } from '@angular/core';
import { AngularFireAuth } from '@angular/fire/compat/auth';
import firebase from 'firebase/compat/app';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  public user:any

  constructor(private authService: AuthService) { }

  ngOnInit(): void {

  }


login() {
    this.authService.login();
    this.user = this.authService.returnUser();
  }
  logout() {
    this.authService.logout();
  }

  showPrediction() {

  }
}
