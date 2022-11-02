import {useMenuStore} from "../core/stores/menuStore";

function App() {
    const showMenu = useMenuStore(e => e.data.show);

    const updateMenu = useMenuStore(state => state.setData)

    return (
        <button onClick={() => updateMenu({show: !showMenu})}>
            Update Menu
        </button>
    )
}

export default App
