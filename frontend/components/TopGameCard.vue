<script setup>
const props = defineProps({
    game: {
        default: SAMPLE_GAME
    }
})
</script>

<template>
    <div class="top-game-card">
        <img class="top-game-thumbnail" :src="game.thumbnail">
        <p class="top-game-title"> {{ game.title }} </p>
        <p class="top-game-description"> {{ game.description }} </p>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';

.top-game-card {
    display: flex;
    flex-direction: column;
    background-color: colors.$light-bg-color;
    border: 4px solid;
    border-color: navy;
}

.top-game-card:hover {
    border-color: crimson;
}

.top-game-title {
    font-size: 18px;
    font-weight: bold;
    margin: 0 0 4px;
}

.top-game-description {
    font-size: 14px;
}

.top-game-card a {
    text-decoration: none;
    color: black;
    border-style: solid;
}

.top-game-thumbnail {
    max-height: 300px;
    max-width: 300px;
}
</style>