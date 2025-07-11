<script setup>
const runConf = useRuntimeConfig()
let currentTab = ref(0)
let engine = ref(null)

let engineBinaryPath = '/thirdparty/godot/godot.editor'
let engineJsPath = '/thirdparty/godot/godot.editor.js'

const persistentPaths = ['/home/web_user'];

let editorCanvas
let gameCanvas
let editorConfig
let gameConfig
let forceGameCanvasReload = ref(0)

let gameName = defineModel({ default: "Новое демо" })

let editor = shallowRef(null);
let game = shallowRef(null);
window.video_driver = '';

function Execute(args) {
    const is_editor = args.filter(function (v) {
        return v === '--editor' || v === '-e';
    }).length !== 0;
    const is_project_manager = args.filter(function (v) {
        return v === '--project-manager';
    }).length !== 0;
    const is_game = !is_editor && !is_project_manager;

    if (video_driver) {
        args.push('--rendering-driver', video_driver);
    }

    if (is_game && (game.value) !== null) {
        alert("A game is already running. Close it first")
        return
    }

    if (is_game) {
        game.value = new Engine(gameConfig)
        game.value.init().then(function () {
            game.value.start({ 'args': args, 'canvas': gameCanvas })
        })
    } else {
        editor.value = new Engine(editorConfig)
        editor.init().then(function () {
            editor.start({ 'args': args, 'canvas': editorCanvas })
        });
    }
}

function closeGame() {
    if (game.value !== null) {
        game.value.requestQuit();
    }
}

onMounted(() => {
    editorCanvas = document.getElementById('editor-canvas');
    gameCanvas = document.getElementById('game-canvas');

    editorConfig = {
        'unloadAfterInit': false,
        'canvasResizePolicy': 0,
        'onExecute': Execute,
        'canvas': editorCanvas,
        // 'onExit': function () {},
        'persistentPaths': persistentPaths
    }

    gameConfig = {
        'persistentPaths': persistentPaths,
        'unloadAfterInit': false,
        'canvas': gameCanvas,
        'canvasResizePolicy': 1,
        'onExit': function () {
            forceGameCanvasReload.value += 1
            nextTick().then(() => {
                gameCanvas = document.getElementById('game-canvas');
                game.value = null
            })
        }
    }

    // Godot has already been loaded by the previus page the user has navigated
    // here from. He moved from play to edit for example. We can safely start
    // Godot
    if (typeof Engine == 'function') {
        loadEngine()
    }
})

let loadEngine = () => {
    editor.value = new Engine(editorConfig)
    editor.init(engineBinaryPath)
        .then(() => {
            const args = ['--project-manager', '--single-window']
            editor.start({ 'args': args, 'persistentDrops': true })
        })
}
 
useHead({
    script: [
        {
            hid: 'godot',
            src: engineJsPath,
            defer: true,
            async: true,
            onload: loadEngine
        }
    ]
})

// Cleanup
onBeforeRouteLeave(() => {
    if (game.value) {
        game.requestQuit();
        game.value = undefined
    }
    if (editor.value) {
        editor.requestQuit()
        editor.value = undefined
    }
})
</script>

<template>
    <Navbar :is-ssr="false"></Navbar>
    <div class="centering-containter">
        <div id="editor-area">
            <input id="game-title-input" class="game-title" type="text" v-model="gameName"></input>
            <div id="tab-navigation" style="display: flex;">
                <button @click="currentTab = 0" class="button">Editor</button>
                <button @click="currentTab = 1"
                    :class="{ 'button': true, 'button-disabled': game === null }">Game</button>
                <button @click="closeGame()" :class="{ 'button': true, 'button-disabled': game === null }">X</button>
                <div id="site-actions">
                    <button>Опубликовать</button>
                    <button>Сохранить</button>
                </div>
            </div>
            <p v-if="!editor.value">Godot загружается</p>
            <canvas id="editor-canvas" width="1200" height="800" :class="{ 'tab-hidden': currentTab !== 0 }"></canvas>
            <p v-if="currentTab === 1 && game === null">The game is not currently running</p>
            <canvas id="game-canvas" width="800" height="600" :class="{ 'tab-hidden': currentTab !== 1 }"
                :key="forceGameCanvasReload"></canvas>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';

.tab-hidden {
    display: none;
}

#site-actions {
    margin-left: 50px;
}

#centering-containter {
    display: flex;
    align-items: center;
    justify-content: center;
}

#game-title-input {
    background-color: colors.$light-highlight-color-darker;
}
</style>