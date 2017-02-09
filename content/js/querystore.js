function QueryStore() {
	riot.observable(this)
	
	var self = this
	
	self.query = {
		preparation:{
			adlib:[],
			follower:[],
		},
		jobs:[
		],
		regularflag:false,
	}
	
	self.on('query_init', ()=>{
		self.trigger('adlib_changed', self.query.preparation.adlib)
		self.trigger('jobs_changed', self.query.jobs)
		self.trigger('follower_changed', self.query.preparation.follower)
	})
	
	self.on('query_export', ()=>{
		self.trigger('query_export_data', self.query)
	})
	
	self.on('query_submit', ()=>{
		self.trigger('query_submited', self.query)
	})
	
	self.on('query_add_adlib', (tag)=>{
		self.query.preparation.adlib.push({list:{listid:"0", tag:tag},userids:[]})
		self.trigger('adlib_changed', self.query.preparation.adlib)
	})
	
	self.on('query_del_adlib',(index)=>{
		self.query.preparation.adlib.splice(index, 1)
		self.trigger('adlib_changed', self.query.preparation.adlib)
	})
	
	self.on('query_add_adlib_user',(index, user_id)=>{
		if(self.query.preparation.adlib[index].userids.indexOf(user_id) === -1){
			self.query.preparation.adlib[index].userids.push(user_id)
			self.trigger('adlib_changed', self.query.preparation.adlib)
		}
	})
	
	self.on('query_del_adlib_user',(index, user_id)=>{
		var userindex=self.query.preparation.adlib[index].userids.indexOf(user_id)
		self.query.preparation.adlib[index].userids.splice(userindex,1)
		self.trigger('adlib_changed', self.query.preparation.adlib)
	})
	
	self.on('query_add_follower', (tag)=>{
		self.query.preparation.follower.push({list:{listid:"0", tag:tag},userid:0})
		self.trigger('follower_changed', self.query.preparation.follower)
	})
	
	self.on('query_del_follower', (index)=>{
		self.query.preparation.follower.splice(index, 1)
		self.trigger('follower_changed', self.query.preparation.follower)
	})
	
	self.on('query_change_follower', (index, user_id)=>{
		self.query.preparation.follower[index].userid = user_id
		self.trigger('follower_changed', self.query.preparation.follower)
	})
	
	self.on('query_add_jobs', ()=>{
		self.query.jobs.push({
			operator:"+",
			listone:{listid:"0",tag:""},
			listanother:{listid:"0",tag:""},
			listresult:{listid:"0",tag:""},
			config:{name:"",publicflag:false,saveflag:false},
		})
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_del_jobs', (index)=>{
		self.query.jobs.splice(index, 1)
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_listone_listid', (index, list_id)=>{
		self.query.jobs[index].listone.listid=list_id
		self.query.jobs[index].listone.tag=""
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_listone_tag', (index, tag)=>{
		self.query.jobs[index].listone.listid="0"
		self.query.jobs[index].listone.tag=tag
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_listanother_listid', (index, list_id)=>{
		self.query.jobs[index].listanother.listid=list_id
		self.query.jobs[index].listanother.tag=""
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_listanother_tag', (index, tag)=>{
		self.query.jobs[index].listanother.listid="0"
		self.query.jobs[index].listanother.tag=tag
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_operator', (index, operator)=>{
		self.query.jobs[index].operator=operator
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_NewSave', (index, name, publicflag)=>{
		self.query.jobs[index].listresult.listid="0"
		self.query.jobs[index].listresult.tag=""
		self.query.jobs[index].config.name=name
		self.query.jobs[index].config.publicflag=publicflag
		self.query.jobs[index].config.saveflag=true
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_UpdateSave', (index, list_id)=>{
		self.query.jobs[index].listresult.listid=list_id
		self.query.jobs[index].listresult.tag=""
		self.query.jobs[index].config.name=""
		self.query.jobs[index].config.publicflag=false
		self.query.jobs[index].config.saveflag=true
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	self.on('query_change_jobs_job_NotSave', (index, tag)=>{
		self.query.jobs[index].listresult.listid="0"
		self.query.jobs[index].listresult.tag=tag
		self.query.jobs[index].config.name=""
		self.query.jobs[index].config.publicflag=false
		self.query.jobs[index].config.saveflag=false
		self.trigger('jobs_changed', self.query.jobs)
	})
	
	
	
	self.userIdscreenNameMap={}
	
	self.on('userIdscreenNameMap_change', (user_id, screen_name)=>{
		self.userIdscreenNameMap[user_id]=screen_name
		self.trigger('userIdscreenNameMap_changed', self.userIdscreenNameMap)
	})
	
	
	self.listIdNameMap={}
	
	self.on('listIdNameMap_change', (list_id, Name)=>{
		self.listIdNameMap[list_id]=Name
		self.trigger('listIdNameMap_changed', self.listIdNameMap)
	})
}