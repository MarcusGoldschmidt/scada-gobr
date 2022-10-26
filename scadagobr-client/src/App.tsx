import AppRouter from "./infra/components/AppRouter";
import {ReactQueryProvider} from "./infra/react-query";

function App() {
    return <>
        <ReactQueryProvider>
            <AppRouter></AppRouter>
        </ReactQueryProvider>
    </>
}

export default App
