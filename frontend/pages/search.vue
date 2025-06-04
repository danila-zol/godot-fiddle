<script setup>
let searchQuery = defineModel()
let results = useState("search-results", () => [])
let maxDescriptionLen = 200

function onSearch() {
    results.value = Array(8).fill(SAMPLE_GAME2)
    results.value = results.value.map((demo) => {
        if (demo.description.length < maxDescriptionLen) {
            return demo
        } else {
            let demoCropped = Object.create(demo)
            demoCropped.description = demoCropped.description.slice(0, maxDescriptionLen) + " ..."
            return demoCropped
        }
    })
}

onUnmounted(() => results.value = [])
</script>

<template>
    <Navbar></Navbar>
    <div class="centering-container">
        <div id="search-area">
            <form @submit.stop.prevent="onSearch" class="search-field">
                <input v-model="searchQuery" type="text" aria-label="search"></input>
            </form>
            <SearchGameCard v-for="result in results" :game="result" :thumbnail-height="'100px'" :thumbnail-width="'100px'" class="search-results"></SearchGameCard>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';

#search-area {
    background-color: colors.$base-color;
    display: flex;
    width: 500px;
    min-height: 100vh;
    flex-direction: column;
    align-items: center;
}

.search-field::before {
    content: "🔎"
}

.search-field {
    width: 300px;
    height: 40px;
    font-size: 24px;
    input {
        height:80%;
        width: 80%;
        background-color: colors.$light-highlight-color-darker;
        border-color: black;
    }
}

.centering-container {
    display: flex;
    align-items: center;
    justify-content: center;
}

.search-results {
    width: 90%;
    margin: 3px;
}
</style>