<template>
    <v-row>
        {{reservedDaysLoop}}
        <v-col>
            <v-row v-for="day in worksDays" :key="day.id" v-show="!day.reserved" align="center" color="text-danger">
                <v-radio-group v-model="time" :mandatory="false" row style="margin:0">
                    <v-radio color="white" :label="day.time" :value="day.time"></v-radio>
                </v-radio-group>
            </v-row>
        </v-col>
        {{numChange}}
    </v-row>
</template>

<script>
    export default {
        name: "picker",
        props: [
            'reservedDays',
            'date',
        ],
        data() {
            return{
                time:'',
                worksDays:
                    [
                    {time:'08:00:00', reserved:false},
                    {time:'08:30:00', reserved:false},
                    {time:'09:00:00', reserved:false},
                    {time:'09:30:00', reserved:false},
                    {time:'10:00:00', reserved:false},
                    {time:'10:30:00', reserved:false},
                    {time:'11:00:00', reserved:false},
                    {time:'11:30:00', reserved:false},
                    {time:'12:00:00', reserved:false},
                    {time:'12:30:00', reserved:false},
                    {time:'13:00:00', reserved:false},
                    {time:'13:30:00', reserved:false},
                    {time:'14:00:00', reserved:false},
                    {time:'14:30:00', reserved:false},
                    {time:'15:00:00', reserved:false},
                    {time:'15:30:00', reserved:false},
                    ],
            }
        },
        computed:{
            numChange(){
                this.$emit('onNumChange', this.time);
            },
            reservedDaysLoop(){
                console.log(this.reservedDays)
                this.worksDays.forEach((item,po) => {
                    for (var index = 0; index < this.reservedDays.length ; index++){
                        if (this.date == this.reservedDays[index].day){
                            if (item.time == this.reservedDays[index].from){
                                this.worksDays[po].reserved = true
                                break;
                            }
                            this.worksDays[po].reserved = false
                        }
                        this.worksDays[po].reserved = false
                    }
                })
            },
        }
    }
</script>
<style>
    .v-label {
        color: red;
    }
</style>