import "./App.css";
import {
  MessageProvider,
  useMessageQuery,
} from "./graphql/generated/codegen.generated";

function App() {
  const { data, loading, error } = useMessageQuery({
    variables: {
      id: `1`,
      provider: MessageProvider.Postgres,
    },
  });

  let content;
  if (loading) {
    content = <div>Loading...</div>;
  } else if (error) {
    content = <div>Error {error.message}</div>;
  } else {
    content = <div>{data?.message.message.data}</div>;
  }

  return (
    <div>
      <h2>Micro Backends 🚀</h2>
      <div>{content}</div>
      <div>{import.meta.env.DEV}</div>
    </div>
  );
}

export default App;
