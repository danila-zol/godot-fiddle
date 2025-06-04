<script setup>
const props = defineProps({
    game: {
        default: SAMPLE_GAME
    },
    gameId: { default: 1 }
})
let maxDescriptionLen = 400
let gameDescriptionCropped = computed(() => {
    let gameDesc = props.game.description
    if (gameDesc.length < maxDescriptionLen) {
        return gameDesc
    } else {
        return gameDesc.slice(0, maxDescriptionLen) + " ..."
    }
})
</script>

<template>
    <NuxtLink :to="'/games/' + gameId" class="top-game-card">
        <img class="top-game-thumbnail" alt="Картинка игры" width="300px" height="300px" :src="game.thumbnail"
            fetchpriority=high>
        <p class="top-game-title"> {{ game.title }} </p>
        <p class="top-game-description"> {{ gameDescriptionCropped }} </p>
    </NuxtLink>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';

.top-game-card {
    display: flex;
    flex-direction: column;
    background-color: colors.$base-color;
    border: 4px solid;
    border-color: colors.$light-highlight-color-darkererer;
}

.top-game-card:hover {
    border-color: colors.$light-highlight-color-darkerer;
}

.top-game-title {
    font-size: 18px;
    font-weight: bold;
    margin: 0 0 4px;
}

.top-game-description {
    font-size: 14px;
}

.top-game-card {
    text-decoration: none;
    color: black;
    border-style: solid;
}
</style>