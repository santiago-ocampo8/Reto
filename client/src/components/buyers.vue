<template>

  
    <v-simple-table>
    
    <template v-slot:default>
       <v-text-field
            v-model="searchBuyer"
            label="Search Buyer by id"
        
          ></v-text-field>
      <thead>
        <tr>
          <th class="text-left">
            Id
          </th>
          <th class="text-left">
            Name
          </th>
          <th class="text-left">
            Age
          </th>
          <th class="text-left">
            Buys
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="item in searchBuyers"
          :key="item.name"
        >
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>{{ item.age }}</td>
          <td><v-btn block color="red darken-4" dark @click="infoBuyer(item.id)">Details</v-btn></td>
     
        </tr>
      </tbody>
    </template>
  </v-simple-table>

</template>

<script>
export default {
  name: 'buyers',
  props: {
    msg: String
  },
  data:()=>({
    Buyers:[],
    searchBuyer:""
    }),
  computed:{
    searchBuyers(){
      return this.Buyers.filter(item=>{
        return item.id.includes(this.searchBuyer)
      })
    }
  },
  mounted(){
    this.getBuyers();
  
  },
  methods:{
    getBuyers(){
      const url="http://localhost:3000/"+"buyers";
      this.axios.get(url).then(response=>{
        var temporal=JSON.stringify(response.data.find[0])
        var temporal2=temporal.replace("@groupby","group");
  
        var final=JSON.parse(temporal2)
        console.log(final)
        this.Buyers=final.group
        
      })
    },
    infoBuyer(id){
      this.$router.push({name:'DetailsBuyer',params:{id:id}})
    }
  }
  

        
}


</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">

</style>
