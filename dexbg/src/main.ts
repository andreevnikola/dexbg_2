import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app/app.module';
import { environment } from './environments/environment';

if (environment.production) {
  enableProdMode();
}

async function buildDex(){
  const key = localStorage.getItem("key");
  if (key && !sessionStorage.getItem("username")) {
    await fetch(`http://localhost:3000/users/authenticate/${key}`)
    .then((res) => {
      if (res.status === 500) {
        alert("Something went wrong with authentication!");
        return
      }
      if (res.status === 401) {
        localStorage.removeItem("key");
        window.location.href = "/auth/signin"
      }

      return res.json()
    })
    .then((data) => {
      sessionStorage.setItem("id", data.id);
      sessionStorage.setItem("username", data.username);
      sessionStorage.setItem("mail", data.mail);
      sessionStorage.setItem("phone", data.phone);
      sessionStorage.setItem("fullname", data.fullname);
      sessionStorage.setItem("gender", data.gender);
    })
  }

  platformBrowserDynamic().bootstrapModule(AppModule)
    .catch(err => console.error(err));
}

buildDex();