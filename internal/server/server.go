package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-embargo/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/flags",s.list);s.mux.HandleFunc("POST /api/flags",s.create);s.mux.HandleFunc("GET /api/flags/{id}",s.get);s.mux.HandleFunc("PUT /api/flags/{id}",s.update);s.mux.HandleFunc("DELETE /api/flags/{id}",s.del)
s.mux.HandleFunc("POST /api/flags/{id}/toggle",s.toggle)
s.mux.HandleFunc("GET /api/evaluate",s.evaluate);s.mux.HandleFunc("POST /api/evaluate",s.evaluateBatch)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/embargo/"})})
return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"flags":oe(s.db.List())})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var f store.Flag;json.NewDecoder(r.Body).Decode(&f);if f.Key==""{we(w,400,"key required");return};if f.Name==""{f.Name=f.Key};if f.Percentage<=0{f.Percentage=100};s.db.Create(&f);wj(w,201,s.db.Get(f.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){f:=s.db.Get(r.PathValue("id"));if f==nil{we(w,404,"not found");return};wj(w,200,f)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.Get(id);if ex==nil{we(w,404,"not found");return};var f store.Flag;json.NewDecoder(r.Body).Decode(&f);if f.Name==""{f.Name=ex.Name};if f.Percentage<=0{f.Percentage=ex.Percentage};s.db.Update(id,&f);wj(w,200,s.db.Get(id))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)toggle(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");s.db.Toggle(id);wj(w,200,s.db.Get(id))}
func(s *Server)evaluate(w http.ResponseWriter,r *http.Request){key:=r.URL.Query().Get("key");enabled,f:=s.db.Evaluate(key);wj(w,200,map[string]any{"key":key,"enabled":enabled,"flag":f})}
func(s *Server)evaluateBatch(w http.ResponseWriter,r *http.Request){var req struct{Keys []string `json:"keys"`};json.NewDecoder(r.Body).Decode(&req);result:=map[string]bool{};for _,k:=range req.Keys{e,_:=s.db.Evaluate(k);result[k]=e};wj(w,200,result)}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"total":s.db.Count(),"enabled":s.db.EnabledCount()})}
func(s *Server)health(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"status":"ok","service":"embargo","flags":s.db.Count(),"enabled":s.db.EnabledCount()})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
