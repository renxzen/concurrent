import { Injectable } from '@angular/core';
import { AngularFireAuth } from '@angular/fire/compat/auth';
import { Router } from '@angular/router';
import firebase from 'firebase/compat/app';
import { Observable } from 'rxjs';
import { CreateUser, User } from '../models/user';
import { HttpClient, HttpParams } from '@angular/common/http';
import { CResponse } from '../models/response';
import { Datagrid, Login, Token } from '../models/transfer';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private user!:any;

  public host:string = "http://localhost:8080";
  // private jwtHelper = new JwtHelperService();
  private token:string = "";
  private username:string = "";

  constructor(
    public auth: AngularFireAuth,
    private router: Router,
    private http:HttpClient
  ) {
    this.auth.onAuthStateChanged((user) => {
      if (user != null) {
        this.router.navigate(['/dashboard'])
        this.user = user
      }
      if (user == null) {
        this.router.navigate(['/'])
      }
    })
   }

  login() {
    this.auth.signInWithPopup(new firebase.auth.GoogleAuthProvider());
  }
  logout() {
    this.auth.signOut();
  }
  returnUser() {
    return this.user
  }

  public register(createUser:CreateUser):Observable<CResponse> {
    return this.http.post<CResponse>(`${this.host}/register`, createUser);
  }

  public loginC(loginUser:Login):Observable<Token> {
    return this.http.post<Token>(`${this.host}/login`, loginUser);
  }

  public getUser():Observable<User> {
    this.loadInfo();
    const params:HttpParams = new HttpParams().set("email", this.username);
		return this.http.get<User>(`${this.host}/find`, { params : params });
  }

  public prediction(datagrid:Datagrid):Observable<CResponse> {
    // this.loadInfo();
    // const params:HttpParams = new HttpParams().set("email", this.username);
    // return this.http.post<CResponse>(`${this.host}/predict`, datagrid, { params : params });
    return this.http.post<CResponse>(`${this.host}/predict`, datagrid);
  }

  public saveInfo(token:string, user:User):void {
    this.token = token;
    localStorage.setItem('token', token);

    this.username = user.email;
    localStorage.setItem('username', user.email);
  }

  public loadInfo():void {
    this.token = localStorage.getItem('token') || "";
    this.username = localStorage.getItem('username') || "";
  }

  public getToken():string {
    return this.token;
  }

  public getUsername():string {
    return this.username;
  }

  public logOut():void {
    this.token = "";
    this.username = "";
    localStorage.removeItem('token');
    localStorage.removeItem('username');
  }

  // public getLoginStatus():boolean {
  //   this.loadInfo();

  //   if (this.token !== ""){
  //     if (this.jwtHelper.decodeToken(this.token).sub != undefined || "") {
  //       if (!this.jwtHelper.isTokenExpired(this.token)) {
  //         if (this.username == this.jwtHelper.decodeToken(this.token).sub) {
  //           return true
  //         }
  //       }
  //     }
  //   }

  //   this.logOut();
  //   return false;
  // }


}
