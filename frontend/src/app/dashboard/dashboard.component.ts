import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { User } from 'firebase/auth';
import { CResponse } from '../core/models/response';
import { Datagrid } from '../core/models/transfer';
import { AuthService } from '../core/services/auth.service';

@Component({
  selector:'app-dashboard',
  templateUrl:'./dashboard.component.html',
  styleUrls:['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  public user!:User;
  public prediction:boolean = true;
  public show:boolean = false;
  public age:number = 0;
  public gender:string = "Male"
  public nDosis:number = 0;
  public data:Datagrid = {} as Datagrid;
  public fabricante:string = "";

  constructor(
    private authService:AuthService
  ) { }

  public ngOnInit():void {
    this.user = this.authService.returnUser();
    console.log(this.user);
  }

  public logout() {
    this.authService.logout();
  }

  public getPrediction():void {
    this.show = false;
    this.data = {
      age:this.age,
      gender:this.gender,
      firstVaccine:this.nDosis == 0? "Ninguna" :this.fabricante,
      secondVaccine:this.nDosis == 2? this.fabricante :"Ninguna",
    }

    this.authService.prediction(this.data)
      .subscribe(
        (res:CResponse) => {
          this.prediction = res.message == "Alive" ? true :false
          this.show = true;
          console.log(res);
      }, (err:HttpErrorResponse) => {
        console.log(err);
      }
    );

    console.log(this.data);
  }


}
