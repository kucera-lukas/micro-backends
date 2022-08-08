import { MessageProvidersProvider } from "./context/providers.context";
import Router from "./router/router";

const App = (): JSX.Element => {
  return (
    <MessageProvidersProvider>
      <Router />
    </MessageProvidersProvider>
  );
};

export default App;
