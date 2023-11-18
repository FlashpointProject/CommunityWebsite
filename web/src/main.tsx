import { createRoot } from 'react-dom/client';
import { App } from './App';
import { CookiesProvider, withCookies } from 'react-cookie';
import { Provider } from 'react-redux';
import store from './store';

import 'semantic-ui-less/semantic.less';
import 'react-calendar/dist/Calendar.css';
import './semantic-ui/core.css';

import { BrowserRouter } from 'react-router-dom';

(async () => {
  const root = document.getElementById('app');

  const ConnectedApp = withCookies(App);

  if (root) {
    createRoot(root).render(
      <Provider store={store}>
        <CookiesProvider>
          <BrowserRouter>
            <ConnectedApp/>
          </BrowserRouter>
        </CookiesProvider>
      </Provider>
    );
  }
})();
