<main>
	<form onsubmit={ export } >
		<button>export</button>
	</form>
	<adlibs />
	<followers />
	<jobs />
	<submit />
	
	
	<script>
		var self=this
	
		self.on('mount', ()=>{
			RiotControl.trigger('query_init')
		})

		RiotControl.on('query_export_data', (query)=>{
			console.log(JSON.stringify(query))
		})

		export(e) {
			RiotControl.trigger('query_export')
			return e.preventDefault()
		}

		adlibpatch(input) {
			query.preparation.adlib=input;
			this.update()
		}
		
	</script>
</main>

<adlibs>
	<h1>adlib</h1>
	<form onsubmit={ add } >
		tag:<input>
		<button>Add</button>
	</form>
	
	<li each={ v,i in adlibs }>
		<adlib data={v} index={i} />
	</li>

	<script>
		var self = this
		self.disabled = true

		self.adlibs=[]
		
		RiotControl.on('adlib_changed', (adlibs)=>{
			self.adlibs = adlibs
			self.update()
		})

		add(e) {
			RiotControl.trigger('query_add_adlib', e.target[0].value)
			return e.preventDefault();
		}

		
	</script>
</adlibs>

<adlib>
	<a onclick={ remove } class="glyphicon glyphicon-remove"></a>
	{opts.data.list.tag}<br />
	<span each={ i in opts.data.userids }>
		<a value={i} onclick={ delmember } class="glyphicon glyphicon-remove"></a>
		{showscreenname(i)}
	</span>
	<br />
	screenname:@<input onBlur={fetch_user_lists} />
	<form style="display: inline" onsubmit={ addmember }>
		<select id="select" class="list-select" >
			<option each={selectlist} value={this[1]}>{this[0]}</option>
		</select>
		<button>Add</button>
	</form>


	<script>
		var self = this
	
		selectlist=[];


		remove() {
			RiotControl.trigger('query_del_adlib', opts.index)
		}

		addmember(e) {
			if(e.target[0].value){
				RiotControl.trigger('query_add_adlib_user', opts.index, e.target[0].value)
			}
			return 	e.preventDefault();
		}

		delmember(e) {
			RiotControl.trigger('query_del_adlib_user', opts.index, e.srcElement.value)
		}

		fetch_user_lists(obj){
			screen_name = obj.srcElement.value;
			$.post("/api/searchuser",{username:screen_name},function(result){
				if(result.status==="ok"){
					this.selectlist = result.data;
					for(i in result.data) {
						RiotControl.trigger('userIdscreenNameMap_change', result.data[i][1], result.data[i][0])
					}
					this.update();
				}
			}.bind(this));
		}.bind(this);
		
		self.screenNameIdMap={}

		RiotControl.on('userIdscreenNameMap_changed', (map)=>{
			self.screenNameIdMap = map
			self.update()
		})

		showscreenname(userID) {
			if (userID in self.screenNameIdMap ) {
				return self.screenNameIdMap[userID];
			} else {
				screenNameIdMap[userID]="";
				$.post("/api/getusers",{userids:userID},function(result){
					if(result.status==="ok"){
						RiotControl.trigger('userIdscreenNameMap_change', result.data[0][1], result.data[0][0])
					} else {
						delete screenNameIdMap[userID];
					}
				}).done(this.update).fail(()=>{
					delete screenNameIdMap[userID];
				})
			}
		}
	</script>
</adlib>

<followers>
	<h1>Follower</h1>
	<form onsubmit={ add } >
		tag:<input>
		<button>Add</button>
	</form>
	
	<li each={ v,i in follower }>
		<follower data={v} index={i} />
	</li>

	<script>
		var self = this
		self.disabled = true

		self.follower=[]
		
		RiotControl.on('follower_changed', (follower)=>{
			self.follower = follower
			self.update()
		})

		add(e) {
			RiotControl.trigger('query_add_follower', e.target[0].value)
			return e.preventDefault();
		}

		
	</script>
</followers>

<follower>
	<a onclick={ remove } class="glyphicon glyphicon-remove"></a>
	{opts.data.list.tag}<br />
		{this.showscreenname(opts.data.userid)}
	<br />
	screenname:@<input onBlur={fetch_user_lists} />
	<form style="display: inline" onsubmit={ change }>
		<select id="select" class="list-select" >
			<option each={selectlist} value={this[1]}>{this[0]}</option>
		</select>
		<button>Change</button>
	</form>


	<script>
		var self = this
	
		selectlist=[];


		remove() {
			RiotControl.trigger('query_del_follower', opts.index)
		}

		change(e) {
			
		console.log(opts.data)
			if(e.target[0].value){
				RiotControl.trigger('query_change_follower', opts.index, e.target[0].value)
			}
			return 	e.preventDefault();
		}

		fetch_user_lists(obj){
			screen_name = obj.srcElement.value;
			$.post("/api/searchuser",{username:screen_name},function(result){
				if(result.status==="ok"){
					this.selectlist = result.data;
					for(i in result.data) {
						RiotControl.trigger('userIdscreenNameMap_change', result.data[i][1], result.data[i][0])
					}
					this.update();
				}
			}.bind(this));
		}.bind(this);
		
		self.screenNameIdMap={}

		RiotControl.on('userIdscreenNameMap_changed', (map)=>{
			self.screenNameIdMap = map
			self.update()
		})

		showscreenname(userID) {
			if (userID===0)return "";
			if (userID in self.screenNameIdMap ) {
				return self.screenNameIdMap[userID];
			} else {
				self.screenNameIdMap[userID]="";
				$.post("/api/getusers",{userids:[userID]},function(result){
					if(result.status==="ok"){
						RiotControl.trigger('userIdscreenNameMap_change', result.data[0][1], result.data[0][0])
					} else {
						delete self.screenNameIdMap[userID];
					}
				}).done(this.update).fail(()=>{
					delete self.screenNameIdMap[userID];
				})
			}
		}
	</script>
</follower>


<jobs>
	<h1>Job</h1>
	<form onsubmit={ add } >
		<button>Add</button>
	</form>
	
	<li each={v,i in this.jobs }>
		<job data={v} index={i} />
	</li>
	
	<script>
		var self = this
	
		self.listidname
		
		self.jobs


		
		RiotControl.on('jobs_changed', (jobs)=>{
			self.jobs = jobs
			self.update()
		})
	
		add(e){
			RiotControl.trigger('query_add_jobs')
			return e.preventDefault();
		}
		
		registerlistname(listID, listName){
			listidname[listID]=listName;
		}
	</script>
</jobs>

<job>
	<a onclick={ remove } class="glyphicon glyphicon-remove"></a><br />
	<joblistselect data={ opts.data.listone } index={ opts.index } oneanother={"one"} />
	<form onchange={ operatorchange }>
	<select id="operator" >
		<option value="+" selected={opts.data.operator==="+"}>+</option>
		<option value="*" selected={opts.data.operator==="*"}>*</option>
		<option value="-" selected={opts.data.operator==="-"}>-</option>
	</select>
	</form>
	<joblistselect data={ opts.data.listanother } index={ opts.index } oneanother={"another"} />
	<br />
	=
	<br />
	<joblistresult config={ opts.data.config } list={opts.data.listresult} index={ opts.index } />
	
	<script>
		remove(index) {
			RiotControl.trigger('query_del_jobs', opts.index)
		}
		
		operatorchange(e) {
			RiotControl.trigger('query_change_jobs_job_operator', opts.index, e.target.value)
		}
		
	</script>
</job>

<joblistselect>
	<div style="display:inline-flex">
	<form onchange={ this.change }>
	<select id="select" >
		<option value="Tag" selected>Tag</option>
		<option value="ListID">List</option>
	</select>
	</form>
	
	<input if={this.type === "Tag"} onchange={ tagnamechange } name="test" value={opts.data.tag}>

	<selectlist if={this.type === "ListID"} index={ opts.index } oneanother={ opts.oneanother }/>
	
	</div>
	<script>
		var self = this
	
		self.type = "Tag"
		self.selectlist=[]
		self.listlist=[]
	
		change(e) {
			this.type = e.target.value
			this.update()
		}
		
		tagnamechange(e) {
			if(opts.oneanother==="one") {
				RiotControl.trigger('query_change_jobs_job_listone_tag', opts.index, e.target.value)
			} else if (opts.oneanother==="another") {
				RiotControl.trigger('query_change_jobs_job_listanother_tag', opts.index, e.target.value)
			}
		}
	</script>
</joblistselect>

<joblistresult>
	<div style="display:inline-flex">
		<form onchange={ this.change }>
			<select id="select" >
				<option value="NewSave" selected>NewSave</option>
				<option value="UpdateSave">UpdateSave</option>
				<option value="NotSave">NotSave</option>
			</select>
		</form>

		<div if={this.type==="NewSave"} style="display:inline-flex">
			Name:<input name="input" onchange={NewSaveInput}>
			<form onchange={NewSave}>
				<select name="PrivatePubric">
					<option value="Private">Private</option>
					<option value="Public">Public</option>	
				</select>
			</form>
		</div>
	
	<!-- loginしているユーザーリストに変えるべき -->
		<selectlist if={this.type==="UpdateSave"} changelistid={UpdateSave} />
	
	<div if={this.type==="NotSave"}>
		Tag:<input  onchange={NotSave}>
	</div>
	</div>
	
	<script>
		var self = this
	
		this.type="NewSave";
		change(e) {
			this.type = e.target.value;
			this.update();
		}
		
		
		self.input = ""
		self.PrivatePubric = false
		NewSaveInput(e) {
			self.input = e.srcElement.value
			RiotControl.trigger('query_change_jobs_job_NewSave', opts.index, self.input, self.PrivatePubric)
		}
		NewSave(e) {
			self.PrivatePubric = (e.srcElement.value==="Public")
			RiotControl.trigger('query_change_jobs_job_NewSave', opts.index, self.input, self.PrivatePubric)
		}
		
		
		UpdateSave(id) {
			RiotControl.trigger('query_change_jobs_job_UpdateSave', opts.index, id)
		}
		
		NotSave(e) {
			RiotControl.trigger('query_change_jobs_job_NotSave', opts.index, e.target.value)
		}
		
	</script>
</joblistresult>

<selectlist>
	screenname:@<input id="screenname" onBlur={fetch_user_lists}>
	<select id="select" class="list-select" onBlur={fetch_lists}>
		<option each={selectlist} value={this[1]}>{this[0]}</option>
	</select>
	<select id="select" class="list-select" onBlur={changelistid}>
		<option each={listlist} value={this[1]}>{this[0]}</option>
	</select>
	
	<script>
		fetch_user_lists(obj){
			screen_name = obj.target.value;
			$.post("/api/searchuser",{username:screen_name},function(result){
				if(result.status==="ok"){
					this.selectlist = result.data;
					for(i in result.data) {
						RiotControl.trigger('userIdscreenNameMap_change', result.data[i][1], result.data[i][0])
					}
					this.update();
				}
			}.bind(this));
		}.bind(this);
		
		fetch_lists(obj) {
			userid = obj.target.value;
			$.post("/api/userlist",{userid:userid},function(result){
				if(result.status==="ok"){
					this.listlist = result.data;
					for(i in result.data) {
						RiotControl.trigger('listIdNameMap_change', result.data[i][1], result.data[i][0])
					}
					this.update();
				}
			}.bind(this));
			console.log(this.listlist)
		}.bind(this);
		
		changelistid(obj) {
			if(opts.oneanother==="one"){
				RiotControl.trigger('query_change_jobs_job_listone_listid', opts.index, obj.target.value)
			} else if(opts.oneanother==="another"){
				RiotControl.trigger('query_change_jobs_job_listanother_listid', opts.index, obj.target.value)
			} else {
				opts.changelistid(obj.target.value)
			}
		}
	</script>
</selectlist>

<submit>
	<form onsubmit={ submit } >
		<button>submit</button>
	</form>
	
	<h1>{ status }</h1>
	
	<script>
		var self = this
	
		status=""
		
		RiotControl.on('query_submited', (query)=>{
			$.post("/api/query",{query:JSON.stringify(query)},function(result){
				console.log(status)
				status=result.status
				this.update()
			}.bind(this)).fail(()=>{
				status="fail"
				this.update()
			})
		})
		submit(e){
			RiotControl.trigger('query_submit')
			return e.preventDefault()
		}
	</script>
</submit>
