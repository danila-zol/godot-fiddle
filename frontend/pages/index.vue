<script setup>
import NewGameCard from '~/components/NewGameCard.vue';

const runConf = useRuntimeConfig()
let topDemosUrl = runConf.public.apiRoot + runConf.public.apiDemosPrefix + "?t=1"
let newDemosUrl = runConf.public.apiRoot + runConf.public.apiDemosPrefix + "?t=1"

let { data: topDemos, error: tderr} = await useFetch(topDemosUrl)
let { data: newDemos, error: nderr } = await useFetch(newDemosUrl)

if (tderr.value) {
    topDemos = Array(3).fill(SAMPLE_GAME)
}

if (nderr.value) {
    newDemos = Array(8).fill(SAMPLE_GAME)
}

</script>

<template>
    <Navbar />
    <div id="main-grid">
        <p id="top-games-label">ТОП-3</p>
        <p id="newest-games-label">Новинки</p>
        <TopGameCard v-for="demo in topDemos" :game="demo"></TopGameCard>
        <div id="new-games-panel">
            <NewGameCard v-for="demo in newDemos" :game="demo"></NewGameCard>
        </div>
        <p id="tags-panel-label">Теги</p>
        <TagsDisplay></TagsDisplay>
    </div>
</template>

<style lang="scss" scoped>
#main-grid {
    display: grid;
    gap: 7px;
    grid-template-columns: 1fr 1fr 1fr 2fr;
    grid-template-rows: auto 2fr auto 1fr;
}

@mixin section-label($fsize: 30px) {
    align-self: end;
    justify-self: start;
    margin: 0;
    padding: 0;
    font-size: $fsize;
}

#top-games-label {
    @include section-label();
    grid-column: span 3;
}

#newest-games-label {
    @include section-label()
}

#newest-games-display {
    grid-column: span 1;
    grid-row: span 3
}

#tags-panel-label {
    @include section-label(24px);
    grid-column: span 3
}

#new-games-panel {
    display: grid;
    justify-content: start;
    align-content: start;
    grid-template-columns: auto auto;
    gap: 2px;
}

.tags-panel {
    grid-column: span 3
}
</style>