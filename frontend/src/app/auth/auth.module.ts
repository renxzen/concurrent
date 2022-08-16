import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AuthRoutingModule } from './auth-routing.module';
import { LoginComponent } from './login/login.component';
import { useDeviceLanguage } from '@angular/fire/auth';

@NgModule({
  declarations: [
    LoginComponent
  ],
  imports: [
    CommonModule,
    AuthRoutingModule
  ],
  exports: [
    LoginComponent
  ],
  providers: [
    { provide: useDeviceLanguage, useValue: true },
  ]
})
export class AuthModule { }
