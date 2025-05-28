<script setup>
const runConf = useRuntimeConfig()
let url = runConf.public.apiRoot + runConf.public.apiDemosPrefix + "?t=1"
let topDemos = await fetch(url)
    .then(
        res => res.json()
    ).then(
        res => JSON.parse(res)
    ).catch(
        err => Array(3).fill(SAMPLE_GAME)
    )
</script>

<template>
    <div id="top-demos">
        <TopGameCard v-for="demo in topDemos" :game="demo"></TopGameCard>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';

#top-demos {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    column-gap: 3px;
    background-color: colors.$dark-separator-color;
}
</style>
