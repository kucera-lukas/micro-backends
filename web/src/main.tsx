import App from "./app";
import GraphqlProvider from "./providers/graphql.provider";

import * as ReactDOM from "react-dom/client";

const root = ReactDOM.createRoot(
  document.querySelector(`#root`) as HTMLElement,
);

root.render(
  <GraphqlProvider>
    <App />
  </GraphqlProvider>,
);
